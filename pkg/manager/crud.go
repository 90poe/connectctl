package manager

import (
	"github.com/pkg/errors"

	"github.com/90poe/connectctl/pkg/client/connect"
)

// GetAllConnectors returns all the connectors in a cluster
func (c *ConnectorManager) GetAllConnectors() ([]*ConnectorWithState, error) {
	c.logger.Debug("getting all connectors (with state)")

	existing, err := c.ListConnectors()
	if err != nil {
		return nil, err
	}

	connectors := make([]*ConnectorWithState, len(existing))
	for i, connectorName := range existing {
		connector, err := c.GetConnector(connectorName)
		if err != nil {
			return nil, err
		}
		connectors[i] = connector
	}

	return connectors, nil
}

// GetConnectors returns information about a named connector in the cluster
func (c *ConnectorManager) GetConnector(connectorName string) (*ConnectorWithState, error) {
	connector, resp, err := c.client.GetConnector(connectorName)
	c.logger.WithField("response", resp).Trace("get connector response")
	if err != nil {
		return nil, errors.Wrapf(err, "getting connector %s", connectorName)
	}

	connectorStatus, resp, err := c.client.GetConnectorStatus(connectorName)
	c.logger.WithField("response", resp).Trace("get connector status response")
	if err != nil {
		return nil, errors.Wrapf(err, "getting connector status %s", connectorName)
	}

	withState := &ConnectorWithState{
		Name:           connector.Name,
		ConnectorState: connectorStatus.Connector,
		Config:         connector.Config,
		Tasks:          connectorStatus.Tasks,
	}

	return withState, nil
}

// ListConnectors returns the names of all connectors in the cluster
func (c *ConnectorManager) ListConnectors() ([]string, error) {
	connectors, resp, err := c.client.ListConnectors()
	c.logger.WithField("response", resp).Trace("list connectors response")
	if err != nil {
		return nil, errors.Wrap(err, "getting existing connectors")
	}

	return connectors, nil
}

// Add will add connectors to a cluster
func (c *ConnectorManager) Add(connectors []connect.Connector) error {
	c.logger.Debug("adding connectors")

	for _, connector := range connectors {
		connectLogger := c.logger.WithField("connector", connector.Name)
		connectLogger.Debug("adding connector")

		resp, err := c.client.CreateConnector(connector)
		connectLogger.WithField("response", resp).Trace("create connector response")
		if err != nil {
			return errors.Wrapf(err, "creating connector %s", connector.Name)
		}

		connectLogger.Debug("added connector")
	}

	c.logger.Debug("created connectors")
	return nil
}

// Remove will remove connectors from a cluster
func (c *ConnectorManager) Remove(connectorNames []string) error {
	c.logger.Debug("removing connectors")

	for _, connectorName := range connectorNames {
		connectLogger := c.logger.WithField("connector", connectorName)
		connectLogger.Debug("deleting connector")

		resp, err := c.client.DeleteConnector(connectorName)
		connectLogger.WithField("response", resp).Trace("delete connector response")
		if err != nil {
			return errors.Wrapf(err, "deleting connector %s", connectorName)
		}

		connectLogger.Debug("deleted connector")
	}

	c.logger.Debug("removed connectors")
	return nil
}
