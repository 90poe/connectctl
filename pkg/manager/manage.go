package manager

import (
	"time"

	"github.com/90poe/connectctl/pkg/client/connect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Manage will start the connector manager running and managing connectors
func (c *ConnectorManager) Manage(source ConnectorSource, stopCH <-chan struct{}) error {
	c.logger.Info("running connector manager")

	syncChannel := time.NewTicker(c.config.SyncPeriod).C
	for {
		select {
		case <-syncChannel:
			connectors, err := source()
			if err != nil {
				return errors.Wrap(err, "getting connector configurations")
			}
			if err = c.reconcileConnectors(connectors); err != nil {
				return errors.Wrap(err, "reconciling connectors")
			}
		case <-stopCH:
			c.logger.Info("Shutting down connector manager")
			return nil
		}
	}
}

func (c *ConnectorManager) reconcileConnectors(connectors []*connect.Connector) error {
	c.logger.Debug("reconciling connectors")

	for _, connector := range connectors {
		err := c.reconcileConnector(connector)
		if err != nil {
			return errors.Wrapf(err, "reconciling connector: %s", connector.Name)
		}
	}

	if c.config.AllowPurge {
		err := c.checkAndDeleteUnmanaged(connectors)
		if err != nil {
			return errors.Wrapf(err, "checking for unmanaged connectors to purge")
		}
	}

	if c.config.AutoRestart {
		err := c.autoRestart(connectors)
		if err != nil {
			return errors.Wrap(err, "checking connector status")
		}
	}

	c.logger.Debug("finished reconciling connectors")
	return nil
}

func (c *ConnectorManager) autoRestart(connectors []*connect.Connector) error {
	c.logger.Debug("auto restarting connectors")

	for _, connector := range connectors {
		connectLogger := c.logger.WithField("connector", connector.Name)
		connectLogger.Debug("getting connector status")

		status, resp, err := c.client.GetConnectorStatus(connector.Name)
		c.logger.WithField("response", resp).Trace("get connector status response")

		if err != nil {
			return errors.Wrapf(err, "getting connector status for %s", connector.Name)
		}

		connectLogger.Debugf("connector state is %s", status.Connector.State)

		// Valid statuses are: RUNNING, UNASSIGNED, PAUSED, FAILED
		if status.Connector.State == "FAILED" {
			connectLogger.Info("connector failed, attempting to restart")
			errRestart := c.restartConnector(connector.Name)
			if errRestart != nil {
				connectLogger.Warn("failed to restart connector")
			} else {
				connectLogger.Info("connector restarted")
			}
			continue
		}

		// If the connector isn't failed it could have some tasks that are failed
		for _, taskState := range status.Tasks {
			taskLogger := connectLogger.WithField("taskid", taskState.ID)
			taskLogger.Debugf("task state is %s", taskState.State)
			if taskState.State == "FAILED" {
				taskLogger.Info("task failed, attempting to restart")
				resp, errRestart := c.client.RestartConnectorTask(connector.Name, taskState.ID)
				taskLogger.WithField("response", resp).Trace("restart connector task response")
				if errRestart != nil {
					taskLogger.Warn("failed to restart connector")
				} else {
					taskLogger.Info("restarted connector task")
				}
			}
		}
	}

	c.logger.Debug("finished checking connectors status")
	return nil
}

func (c *ConnectorManager) checkAndDeleteUnmanaged(connectors []*connect.Connector) error {
	c.logger.Debug("purging any unmanaged connectors from cluster")

	existing, resp, err := c.client.ListConnectors()
	c.logger.WithField("response", resp).Trace("list connectors response")
	if err != nil {
		return errors.Wrap(err, "getting existing connectors")
	}

	var unmanaged []string
	for _, existingName := range existing {
		if !containsConnector(existingName, connectors) {
			unmanaged = append(unmanaged, existingName)
		}
	}

	if len(unmanaged) == 0 {
		c.logger.Debug("no unmanaged connectors")
		return nil
	}

	return c.deleteUnmanaged(unmanaged)
}

func (c *ConnectorManager) deleteUnmanaged(unmanagedConnectors []string) error {
	c.logger.Debug("deleting unmanaged connectors")

	err := c.Remove(unmanagedConnectors)
	if err != nil {
		return errors.Wrap(err, "deleting unmanaged connectors")
	}

	c.logger.Debug("deleted unmanaged connectors")

	return nil
}

func (c *ConnectorManager) reconcileConnector(connector *connect.Connector) error {
	connectLogger := c.logger.WithField("connector", connector.Name)
	connectLogger.Debug("reconciling connector")
	defer connectLogger.Debug("reconciled connector")

	connectLogger.Debug("checking if connector exists in cluster")
	existingConnectors, resp, err := c.client.GetConnector(connector.Name)
	connectLogger.WithField("response", resp).Trace("get connector response")

	if err != nil {
		if connect.IsNotFound(err) {
			return c.handleNewConnector(connector, connectLogger)
		}
		return errors.Wrapf(err, "getting existing connector from cluster")
	}

	if existingConnectors != nil {
		return c.handleExistingConnector(connector, existingConnectors, connectLogger)
	}

	connectLogger.Warn("unexpected situation, ignoring connector")
	return nil
}

func (c *ConnectorManager) handleExistingConnector(connector *connect.Connector, existingConnector *connect.Connector, logger *log.Entry) error {
	logger.Debug("connector already exists")

	if existingConnector.ConfigEqual(connector) {
		logger.Debug("connector configuration is the same, no action needed")
		return nil
	}

	logger.Info("connector configuration differs, updating")
	_, resp, err := c.client.UpdateConnectorConfig(existingConnector.Name, connector.Config)
	logger.WithField("response", resp).Trace("update connector config response")

	if err != nil {
		logger.WithError(err).Debugf("error updating connector config: %v", connector.Config)
		return errors.Wrap(err, "updating connector config")
	}

	logger.Info("connector config updated")
	return nil
}

func (c *ConnectorManager) handleNewConnector(connector *connect.Connector, logger *log.Entry) error {
	logger.Info("connector doesn't exist, creating")

	err := c.Add([]*connect.Connector{connector})

	if err != nil {
		logger.WithError(err).Debug("error creating connector")
		return errors.Wrap(err, "creating connector")
	}

	logger.Info("created connector successfully")
	return nil
}

func containsConnector(connectorName string, connectors []*connect.Connector) bool {
	for _, c := range connectors {
		if c.Name == connectorName {
			return true
		}
	}
	return false
}
