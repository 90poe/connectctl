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
	c.logger.Info("restarting specified connectors")

	for _, connectorName := range connectorNames {

		err := c.restartConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "restarting connector %s", connectorName)
		}

		if restartTasks {
			err := c.restartConnectorTasks(connectorName, forceRestartTasks, taskIDs)
			if err != nil {
				return errors.Wrapf(err, "restarting connector %s", connectorName)
			}
		}
	}

	return nil
}

func (c *ConnectorManager) restartConnector(connectorName string) error {
	connectLogger := c.logger.WithField("connector", connectorName)
	connectLogger.Info("restarting connector")

	resp, err := c.client.RestartConnector(connectorName)
	connectLogger.WithField("response", resp).Trace("restart connector response")

	if err != nil {
		return errors.Wrap(err, "calling restart connector API")
	}

	connectLogger.Info("restarted connector")
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

	connectLogger := c.logger.
		WithField("connector", connectorName).
		WithField("taskID", taskIDs)
	connectLogger.Info("restarting connector tasks")

	for _, taskID := range tasks.IDs() {
		err := c.restartConnectorTask(connectorName, taskID)
		if err != nil {
			return err
		}
	}

	connectLogger.Info("restarted connector tasks")
	return nil
}

func (c *ConnectorManager) restartConnectorTask(connectorName string, taskID int) error {
	connectLogger := c.logger.
		WithField("connector", connectorName).
		WithField("taskID", taskID)

	connectLogger.Info("restarting connector task")

	resp, err := c.client.RestartConnectorTask(connectorName, taskID)
	connectLogger.WithField("response", resp).Trace("restart connector response")

	if err != nil {
		return errors.Wrapf(err, "calling restart task connector API for task %d", taskID)
	}

	connectLogger.Info("restarted connector task")
	return nil
}
