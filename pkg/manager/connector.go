package manager

import (
	"fmt"
	"time"

	"github.com/90poe/connectctl/pkg/client/connect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ConnectorSource will return a slice of the desired connector configurationbs
type ConnectorSource func() ([]*connect.Connector, error)

// ConnectorManager manages connectors in a Kafka Connect cluster
type ConnectorManager struct {
	config *Config
	client *connect.Client
	logger *log.Entry
}

// NewConnectorsManager creates a new ConnectorManager
func NewConnectorsManager(config *Config) (*ConnectorManager, error) {
	userAgent := fmt.Sprintf("90poe.io/connectctl/%s", config.Version)

	client, err := connect.NewClient(config.ClusterURL, userAgent)
	if err != nil {
		return nil, errors.Wrap(err, "creating connect client")
	}

	return &ConnectorManager{
		config: config,
		client: client,
		logger: config.Logger,
	}, nil
}

// Run will start the connector manager running and managing connectors
func (c *ConnectorManager) Run(source ConnectorSource, stopCH <-chan struct{}) error {
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

	err := c.checkConnectorStatus(connectors)
	if err != nil {
		return errors.Wrap(err, "checking connector status")
	}

	c.logger.Debug("finished reconciling connectors")
	return nil
}

func (c *ConnectorManager) checkConnectorStatus(connectors []*connect.Connector) error {
	c.logger.Debug("checking connectors status")

	for _, connector := range connectors {
		connectLogger := c.logger.WithField("connector", connector.Name)
		connectLogger.Debug("getting connector status")

		status, resp, err := c.client.GetConnectorStatus(connector.Name)
		c.logger.WithField("response", resp).Trace("get connector status response")

		if err != nil {
			return errors.Wrapf(err, "getting connector status for %s", connector.Name)
		}

		connectLogger.Debugf("connector state is %s", status.Connector.State)

		switch status.Connector.State {
		case "RUNNING":
			break
		case "UNASSIGNED":
			break
		case "PAUSED":
			break
		case "FAILED":
			connectLogger.Error("connector has failed, restarting")
			//TODO: restart
		}

		// TODO loop around the tasks
		for _, taskState := range status.Tasks {
			taskLogger := connectLogger.WithField("taskid", taskState.ID)
			taskLogger.Debugf("task state is %s", taskState.State)
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

	for _, connectorName := range unmanagedConnectors {
		connectLogger := c.logger.WithField("connector", connectorName)
		connectLogger.Info("deleting connector from cluster")

		resp, err := c.client.DeleteConnector(connectorName)
		connectLogger.WithField("response", resp).Trace("delete connector response")

		if err != nil {
			return errors.Wrapf(err, "deleting connector %s", connectorName)
		}

		connectLogger.Info("deleted connector")
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

	//TODO: handle API errors

	logger.Info("connector config updated")
	return nil
}

func (c *ConnectorManager) handleNewConnector(connector *connect.Connector, logger *log.Entry) error {
	logger.Info("connector doesn't exist, creating")
	resp, err := c.client.CreateConnector(connector)
	logger.WithField("response", resp).Trace("create connector response")

	if err != nil {
		logger.WithError(err).Debug("error creating connector")
		return errors.Wrap(err, "creating connector")
	}

	//TODO: handle API error

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
