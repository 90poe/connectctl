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
}

func NewConnectorsManager(clusterURL string, connectors []*connect.Connector, syncPeriod time.Duration) (*ConnectorManager, error) {
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
	}, nil
}

func (c *ConnectorManager) Run(stopCH <-chan struct{}) error {
	log.Infof("running connector manager for server %s\n", c.clusterURL)

	syncChannel := time.NewTicker(c.syncPeriod).C
	for {
		select {
		case <-syncChannel:
			if err := c.reconcileConnectors(); err != nil {
				return errors.Wrap(err, "reconciling connectors")
			}
		case <-stopCH:
			log.Info("Shutting down connector manager")
			return nil
		}
	}

	return nil
}

func (c *ConnectorManager) reconcileConnectors() error {
	log.Debug("reconciling connectors")

	for _, connector := range c.connectors {
		log.WithField("connector", connector.Name).Debug("reconciling connector")

		err := c.reconcileConnector(connector)
		if err != nil {
			return errors.Wrapf(err, "reconciling connector: %s", connector.Name)
		}
	}

	log.Debug("finished reconciling connectors")
	return nil
}

func (c *ConnectorManager) reconcileConnector(connector *connect.Connector) error {
	existingConnectors, resp, err := c.client.GetConnector(connector.Name)
	if err != nil {
		return errors.Wrapf(err, "getting existing connector from cluster")
	}

	if resp.StatusCode == http.StatusNotFound {
		//TODO: create
	}

	if diff := cmp.Diff(connector, existingConnectors); diff != "" {
		log.Infof("Diff detected %s %s", connector.Name, diff)
		// TODO: update the connector
	}

	return nil
}
