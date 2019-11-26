package manager

import "github.com/pkg/errors"

// Restart will restart a number of connectors in a cluster
func (c *ConnectorManager) Restart(connectorNames []string, restartTasks bool,
	forceRestartTasks bool, taskIDs []int) error {
	if len(connectorNames) > 0 {
		return c.restartConnectors(connectorNames, restartTasks, forceRestartTasks, taskIDs)
	}

	connectorNames, err := c.ListConnectors()
	if err != nil {
		return err
	}

	return c.restartConnectors(connectorNames, restartTasks, forceRestartTasks, taskIDs)
}

func (c *ConnectorManager) restartConnectors(connectorNames []string, restartTasks bool,
	forceRestartTasks bool, taskIDs []int) error {
	for _, connectorName := range connectorNames {
		if err := c.restartConnector(connectorName); err != nil {
			return errors.Wrapf(err, "error restarting connector : %s", connectorName)
		}

		if restartTasks {
			if err := c.restartConnectorTasks(connectorName, forceRestartTasks, taskIDs); err != nil {
				return errors.Wrapf(err, "error restarting task : %s", connectorName)
			}
		}
	}

	return nil
}

func (c *ConnectorManager) restartConnector(connectorName string) error {
	if _, err := c.client.RestartConnector(connectorName); err != nil {
		return errors.Wrapf(err, "error calling restart connector : %s", connectorName)
	}

	return nil
}

func (c *ConnectorManager) restartConnectorTasks(connectorName string, forceRestartTasks bool, taskIDs []int) error {
	connector, err := c.GetConnector(connectorName)
	if err != nil {
		return err
	}

	if len(taskIDs) == 0 {
		taskIDs = connector.Tasks.IDs()
	}

	tasks := connector.Tasks.Filter(ByID(taskIDs...))

	if !forceRestartTasks {
		tasks = tasks.Filter(IsNotRunning)
	}

	for _, taskID := range tasks.IDs() {
		if _, err := c.client.RestartConnectorTask(connectorName, taskID); err != nil {
			return errors.Wrapf(err, "error calling restart task connector API for task %d", taskID)
		}
	}

	return nil
}
