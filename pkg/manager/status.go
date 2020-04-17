package manager

import (
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/pkg/errors"
)

// Status - gets status of specified (or all) connectors
func (c *ConnectorManager) Status(connectors []string) ([]*connect.ConnectorStatus, error) {
	if len(connectors) == 0 {
		return c.allConnectorsStatus()
	}

	return c.specifiedConnectorsStatus(connectors)
}

func (c *ConnectorManager) allConnectorsStatus() ([]*connect.ConnectorStatus, error) {
	existing, _, err := c.client.ListConnectors()
	if err != nil {
		return nil, errors.Wrap(err, "error listing connectors")
	}

	return c.specifiedConnectorsStatus(existing)
}

func (c *ConnectorManager) specifiedConnectorsStatus(connectors []string) ([]*connect.ConnectorStatus, error) {
	statusList := make([]*connect.ConnectorStatus, len(connectors))

	for idx, connectorName := range connectors {
		status, _, err := c.client.GetConnectorStatus(connectorName)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting connector status for %s", connectorName)
		}

		statusList[idx] = status
	}

	return statusList, nil
}
