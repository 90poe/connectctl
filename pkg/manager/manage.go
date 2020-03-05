package manager

import (
	"time"

	"github.com/90poe/connectctl/pkg/client/connect"

	"github.com/pkg/errors"
)

// Manage will start the connector manager running and managing connectors
func (c *ConnectorManager) Manage(source ConnectorSource, stopCH <-chan struct{}) error {
	// mark ourselves as having an unhealthy state until we have
	// tried to contact the kafka-connect instance
	c.readinessState = errorState

	syncChannel := time.NewTicker(c.config.SyncPeriod).C
	for {
		select {
		case <-syncChannel:
			err := c.Sync(source)
			if err != nil {
				// set back into an unhealthy state
				c.readinessState = errorState
				return errors.Wrap(err, "error synchronising connectors for source")
			}
			// mark ourselves as being in an ok state as we have
			// started syncing without any error
			c.readinessState = okState
		case <-stopCH:
			return nil
		}
	}
}

// Sync will synchronise the desired and actual state of connectors in a cluster
func (c *ConnectorManager) Sync(source ConnectorSource) error {

	connectors, err := source()
	if err != nil {
		return errors.Wrap(err, "error getting connector configurations")
	}

	// create dynamic restart policy here, overriding with the supplied one (if any)
	policy := map[string]Policy{}

	for _, c := range connectors {
		policy[c.Name] = Policy{
			MaxConnectorRestarts:   0,
			ConnectorRestartPeriod: time.Second * 10,
			MaxTaskRestarts:        0,
			TaskRestartPeriod:      time.Second * 10,
		}
	}

	// apply overrides
	if c.config.RestartPolicy != nil {
		for k, v := range c.config.RestartPolicy.Connectors {

			p := policy[k]

			if v.MaxConnectorRestarts != 0 {
				p.MaxConnectorRestarts = v.MaxConnectorRestarts
			}
			if v.ConnectorRestartPeriod != 0 {
				p.ConnectorRestartPeriod = v.ConnectorRestartPeriod
			}
			if v.MaxTaskRestarts != 0 {
				p.MaxTaskRestarts = v.MaxTaskRestarts
			}
			if v.TaskRestartPeriod != 0 {
				p.TaskRestartPeriod = v.TaskRestartPeriod
			}

			policy[k] = p
		}
	}

	if err = c.reconcileConnectors(connectors, policy); err != nil {
		return errors.Wrap(err, "error synchronising connectors")
	}
	return nil
}

func (c *ConnectorManager) reconcileConnectors(connectors []connect.Connector, restartPolicy map[string]Policy) error {
	for _, connector := range connectors {
		if err := c.reconcileConnector(connector); err != nil {
			return errors.Wrapf(err, "error reconciling connector: %s", connector.Name)
		}
	}

	if c.config.AllowPurge {
		if err := c.checkAndDeleteUnmanaged(connectors); err != nil {
			return errors.Wrapf(err, "error checking for unmanaged connectors to purge")
		}
	}

	if c.config.AutoRestart {
		if err := c.autoRestart(connectors, restartPolicy); err != nil {
			return errors.Wrap(err, "error checking connector status")
		}
	}

	return nil
}

func (c *ConnectorManager) autoRestart(connectors []connect.Connector, policy map[string]Policy) error {

	for _, connector := range connectors {
		status, _, err := c.client.GetConnectorStatus(connector.Name)

		if err != nil {
			return errors.Wrapf(err, "error getting connector status for %s", connector.Name)
		}

		// Valid statuses are: RUNNING, UNASSIGNED, PAUSED, FAILED
		if status.Connector.State == "FAILED" {

			err = c.retryRestartConnector(connector.Name, policy[connector.Name].MaxConnectorRestarts, policy[connector.Name].ConnectorRestartPeriod)

			if err != nil {
				return err
			}
			// we continue here as to give the tasks a chance to restart
			continue
		}

		// If the connector isn't failed it could have some tasks that are failed
		for _, taskState := range status.Tasks {
			if IsNotRunning(taskState) {

				err = c.retryRestartConnectorTask(connector.Name, taskState.ID, policy[connector.Name].MaxTaskRestarts, policy[connector.Name].TaskRestartPeriod)

				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *ConnectorManager) retryRestartConnector(name string, retrys int, retryPeriod time.Duration) error {
	err := c.restartConnector(name)

	if err != nil {
		if retrys == 0 {
			return err
		}
		attempts := 0

		for attempts < retrys {
			time.Sleep(retryPeriod)

			err = c.restartConnector(name)

			if err != nil {
				attempts++
			} else {
				return nil
			}
		}
		return err
	}

	return nil
}

func (c *ConnectorManager) retryRestartConnectorTask(name string, taskID int, retrys int, retryPeriod time.Duration) error {
	_, err := c.client.RestartConnectorTask(name, taskID)

	if err != nil {
		if retrys == 0 {
			return err
		}
		attempts := 0

		for attempts < retrys {
			time.Sleep(retryPeriod)

			_, err = c.client.RestartConnectorTask(name, taskID)

			if err != nil {
				attempts++
			} else {
				return nil
			}
		}
		return err
	}

	return nil
}

func (c *ConnectorManager) checkAndDeleteUnmanaged(connectors []connect.Connector) error {
	existing, _, err := c.client.ListConnectors()
	if err != nil {
		return errors.Wrap(err, "error getting existing connectors")
	}

	var unmanaged []string
	for _, existingName := range existing {
		if !containsConnector(existingName, connectors) {
			unmanaged = append(unmanaged, existingName)
		}
	}

	if len(unmanaged) == 0 {
		return nil
	}

	if err := c.Remove(unmanaged); err != nil {
		return errors.Wrap(err, "error deleting unmanaged connectors")
	}
	return nil
}

func (c *ConnectorManager) reconcileConnector(connector connect.Connector) error {
	existingConnectors, _, err := c.client.GetConnector(connector.Name)

	if err != nil {
		if connect.IsNotFound(err) {
			return c.handleNewConnector(connector)
		}
		return errors.Wrap(err, "error getting existing connector from cluster")
	}

	if existingConnectors != nil {
		return c.handleExistingConnector(connector, existingConnectors)
	}

	return nil
}

func (c *ConnectorManager) handleExistingConnector(connector connect.Connector, existingConnector *connect.Connector) error {
	if existingConnector.ConfigEqual(connector) {
		return nil
	}

	if _, _, err := c.client.UpdateConnectorConfig(existingConnector.Name, connector.Config); err != nil {
		return errors.Wrap(err, "error updating connector config")
	}

	return nil
}

func (c *ConnectorManager) handleNewConnector(connector connect.Connector) error {
	if err := c.Add([]connect.Connector{connector}); err != nil {
		return errors.Wrap(err, "error creating connector")
	}

	return nil
}

func containsConnector(connectorName string, connectors []connect.Connector) bool {
	for _, c := range connectors {
		if c.Name == connectorName {
			return true
		}
	}
	return false
}
