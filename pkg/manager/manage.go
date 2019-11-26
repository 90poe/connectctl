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
	if err = c.reconcileConnectors(connectors); err != nil {
		return errors.Wrap(err, "error synchronising connectors")
	}
	return nil
}

func (c *ConnectorManager) reconcileConnectors(connectors []connect.Connector) error {
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
		if err := c.autoRestart(connectors); err != nil {
			return errors.Wrap(err, "error checking connector status")
		}
	}

	return nil
}

func (c *ConnectorManager) autoRestart(connectors []connect.Connector) error {
	for _, connector := range connectors {
		status, _, err := c.client.GetConnectorStatus(connector.Name)

		if err != nil {
			return errors.Wrapf(err, "error getting connector status for %s", connector.Name)
		}

		// Valid statuses are: RUNNING, UNASSIGNED, PAUSED, FAILED
		if status.Connector.State == "FAILED" {
			// CONSIDER : add feature to track errors restarting connectors
			// Maybe introduce a policy to tolerate a certain amount of attempts
			// then error or fail the healthcheck
			_ = c.restartConnector(connector.Name)
			continue
		}

		// If the connector isn't failed it could have some tasks that are failed
		for _, taskState := range status.Tasks {
			if taskState.State == "FAILED" {
				// CONSIDER : add feature to track errors restarting tasks
				// Maybe introduce a policy to tolerate a certain amount of attempts
				// then error or fail the healthcheck
				_, _ = c.client.RestartConnectorTask(connector.Name, taskState.ID)
			}
		}
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
