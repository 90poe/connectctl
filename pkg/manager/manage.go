package manager

import (
	"fmt"
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
			MaxConnectorRestarts:   1,
			ConnectorRestartPeriod: defaultRestartPeriod,
			MaxTaskRestarts:        1,
			TaskRestartPeriod:      defaultRestartPeriod,
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
		name := connector.Name
		err := c.retryRestartConnector(name, policy[name].MaxConnectorRestarts, policy[name].ConnectorRestartPeriod)
		if err != nil {
			return err
		}
		err = c.retryRestartConnectorTask(name, policy[name].MaxTaskRestarts, policy[name].TaskRestartPeriod)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConnectorManager) retryRestartConnector(name string, retrys int, retryPeriod time.Duration) error {

	attempts := 0
	for attempts <= retrys {
		status, _, err := c.client.GetConnectorStatus(name)

		if err != nil {
			return errors.Wrapf(err, "error getting connector status: %s", name)
		}

		// Valid statuses are: RUNNING, UNASSIGNED, PAUSED, FAILED
		if status.Connector.State == "FAILED" {
			err = c.restartConnector(name)
			if err != nil {
				return errors.Wrapf(err, "error restarting connector: %s", name)
			}

		} else {
			return nil
		}

		attempts++
		time.Sleep(retryPeriod)
	}
	return fmt.Errorf("error restarting connector: %s, retrys: %d", name, retrys)
}

func (c *ConnectorManager) retryRestartConnectorTask(name string, retrys int, retryPeriod time.Duration) error {

	attempts := 0
	for attempts <= retrys {
		status, _, err := c.client.GetConnectorStatus(name)

		if err != nil {
			return errors.Wrapf(err, "error getting connector status: %s", name)
		}

		// Valid statuses are: RUNNING, UNASSIGNED, PAUSED, FAILED
		if status.Connector.State == "RUNNING" {
			runningTasks := 0

			for _, taskState := range status.Tasks {
				if IsNotRunning(taskState) {
					_, err := c.client.RestartConnectorTask(name, taskState.ID)

					if err != nil {
						return err
					}
				} else {
					runningTasks++
				}
			}
			if runningTasks == len(status.Tasks) {
				return nil
			}

		} else {
			return nil
		}

		attempts++
		time.Sleep(retryPeriod)
	}
	return fmt.Errorf("error restarting connector task: %s, retrys: %d", name, retrys)
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
