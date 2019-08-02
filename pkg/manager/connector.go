package manager

import (
	"fmt"
	"net/http"
	"time"

	"github.com/90poe/connectctl/pkg/connect"
	"github.com/90poe/connectctl/pkg/version"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ConnectorManager struct {
	clusterURL string
	connectors []*connect.Connector
	syncPeriod time.Duration
	client     *connect.Client
	logger     *log.Entry
}

func NewConnectorsManager(clusterURL string, connectors []*connect.Connector, syncPeriod time.Duration, logger *log.Entry) (*ConnectorManager, error) {
	userAgent := fmt.Sprintf("90poe.io/connectctl/%s", version.Version)

	client, err := connect.NewClient(clusterURL, userAgent)
	if err != nil {
		return nil, errors.Wrap(err, "creating connect client")
	}

	return &ConnectorManager{
		clusterURL: clusterURL,
		connectors: connectors,
		syncPeriod: syncPeriod,
		client:     client,
		logger:     logger,
	}, nil
}

func (c *ConnectorManager) Run(stopCH <-chan struct{}) error {
	c.logger.Info("running connector manager")

	syncChannel := time.NewTicker(c.syncPeriod).C
	for {
		select {
		case <-syncChannel:
			if err := c.reconcileConnectors(); err != nil {
				return errors.Wrap(err, "reconciling connectors")
			}
		case <-stopCH:
			c.logger.Info("Shutting down connector manager")
			return nil
		}
	}
}

func (c *ConnectorManager) reconcileConnectors() error {
	c.logger.Debug("reconciling connectors")

	for _, connector := range c.connectors {
		err := c.reconcileConnector(connector)
		if err != nil {
			return errors.Wrapf(err, "reconciling connector: %s", connector.Name)
		}
	}

	c.logger.Debug("finished reconciling connectors")
	return nil
}

func (c *ConnectorManager) reconcileConnector(connector *connect.Connector) error {
	connectLogger := c.logger.WithField("connector", connector.Name)
	connectLogger.Info("reconciling connector")

	connectLogger.Debug("checking if connector exists in cluster")
	existingConnectors, resp, err := c.client.GetConnector(connector.Name)
	if err != nil {
		if !connect.IsNotFound(err) {
			return errors.Wrapf(err, "getting existing connector from cluster")
		}
	}

	if resp.StatusCode == http.StatusNotFound {
		connectLogger.Info("connector doesn't exist, creating")
		_, err = c.client.CreateConnector(connector)
		if err != nil {
			connectLogger.WithError(err).Debug("error creating connector")
			return errors.Wrap(err, "creating connector")
		}
		//TODO: handle API error

		connectLogger.Info("created connector sucessfully")
		return nil
	}

	if diff := cmp.Diff(connector, existingConnectors); diff != "" {
		connectLogger.Infof("Diff detected in connector config %s", diff)
		// TODO: update the connector
	}

	connectLogger.Info("reconciled connector")
	return nil
}
