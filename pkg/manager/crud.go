package manager

import (
	"github.com/pkg/errors"

	"github.com/90poe/connectctl/pkg/client/connect"
)

// GetAllConnectors returns all the connectors in a cluster
func (c *ConnectorManager) GetAllConnectors() ([]*ConnectorWithState, error) {

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

// GetConnector returns information about a named connector in the cluster
func (c *ConnectorManager) GetConnector(connectorName string) (*ConnectorWithState, error) {
	connector, _, err := c.client.GetConnector(connectorName)
	if err != nil {
		return nil, errors.Wrapf(err, "getting connector %s", connectorName)
	}

	connectorStatus, _, err := c.client.GetConnectorStatus(connectorName)
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
	connectors, _, err := c.client.ListConnectors()

	if err != nil {
		return nil, errors.Wrap(err, "getting existing connectors")
	}

	return connectors, nil
}

// Add will add connectors to a cluster
func (c *ConnectorManager) Add(connectors []connect.Connector) error {

	for _, connector := range connectors {

		_, err := c.client.CreateConnector(connector)
		if err != nil {
			return errors.Wrapf(err, "error creating connector %s", connector.Name)
		}
	}

	return nil
}

// Remove will remove connectors from a cluster
func (c *ConnectorManager) Remove(connectorNames []string) error {

	for _, connectorName := range connectorNames {

		_, err := c.client.DeleteConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "error deleting connector %s", connectorName)
		}
	}

	return nil
}
