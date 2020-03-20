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
	// successfully configured the kafka-connect instance
	c.readinessState = errorState

	syncChannel := time.NewTicker(c.config.SyncPeriod).C
	for {
		select {
		case <-syncChannel:

			// we only want to try Syncing if we can contact the kafka-connect instance.
			// Using the LivenessCheck as a proxy for calculating the connection
			if c.livenessState == okState {
				err := c.Sync(source)
				if err != nil {
					// set back into an unhealthy state
					c.readinessState = errorState
					return errors.Wrap(err, "error synchronising connectors for source")
				}
				// mark ourselves as being in an ok state as we have
				// started syncing without any error
				c.readinessState = okState
			} else {
				c.logger.Infof("skipping sync as livenessState == %v", c.livenessState)
			}

		case <-stopCH:
			return nil
		}
	}
}

// Sync will synchronise the desired and actual state of connectors in a cluster
func (c *ConnectorManager) Sync(source ConnectorSource) error {
	c.logger.Infof("loading connectors")
	connectors, err := source()
	if err != nil {
		return errors.Wrap(err, "error getting connector configurations")
	}
	c.logger.Infof("connectors loaded : %d", len(connectors))
	// creating a runtime restart policy here, overriding with the supplied one (if any)
	// Ensuring that we have a policy defined for each connector we are manging here
	// dramatically simplifies the management and restart code
	policy := runtimePolicyFromConnectors(connectors, c.config)

	if err = c.reconcileConnectors(connectors, policy); err != nil {
		return errors.Wrap(err, "error synchronising connectors")
	}
	return nil
}

func (c *ConnectorManager) reconcileConnectors(connectors []connect.Connector, restartPolicy runtimeRestartPolicy) error {
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

func (c *ConnectorManager) autoRestart(connectors []connect.Connector, restartPolicy runtimeRestartPolicy) error {
	for _, connector := range connectors {
		name := connector.Name
		err := c.retryRestartConnector(name, restartPolicy[name].ConnectorRestartsMax, restartPolicy[name].ConnectorRestartPeriod)
		if err != nil {
			return err
		}
		err = c.retryRestartConnectorTask(name, restartPolicy[name].TaskRestartsMax, restartPolicy[name].TaskRestartPeriod)
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

		if isConnectorFailed(status.Connector) {
			c.logger.Infof("connector not running: %s", name)

			if err = c.restartConnector(name); err != nil {
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

func isConnectorRunning(c connect.ConnectorState) bool {
	return c.State == "RUNNING" // nolint
}

func isConnectorFailed(c connect.ConnectorState) bool {
	return c.State == "FAILED" // nolint
}

func isTaskFailed(t connect.TaskState) bool {
	return t.State == "FAILED" //nolint
}

func (c *ConnectorManager) retryRestartConnectorTask(name string, retrys int, retryPeriod time.Duration) error {
	attempts := 0
	for attempts <= retrys {
		status, _, err := c.client.GetConnectorStatus(name)

		if err != nil {
			return errors.Wrapf(err, "error getting connector status: %s", name)
		}

		if isConnectorRunning(status.Connector) {
			runningTasks := 0

			for _, taskState := range status.Tasks {
				if isTaskFailed(taskState) {
					c.logger.Infof("task not running: %d ( %s )", taskState.ID, name)

					if _, err := c.client.RestartConnectorTask(name, taskState.ID); err != nil {
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
