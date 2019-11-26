package manager

import "github.com/pkg/errors"

// Pause will pause a number of connectors in a cluster
func (c *ConnectorManager) Pause(connectors []string) error {
	if len(connectors) == 0 {
		return c.pauseAllConnectors()
	}

	return c.pauseSpecifiedConnectors(connectors)
}

func (c *ConnectorManager) pauseAllConnectors() error {

	existing, _, err := c.client.ListConnectors()
	if err != nil {
		return errors.Wrap(err, "error listing connectors")
	}

	for _, connectorName := range existing {
		err := c.pauseConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "pausing connector %s", connectorName)
		}
	}

	return nil
}

func (c *ConnectorManager) pauseSpecifiedConnectors(connectors []string) error {

	for _, connectorName := range connectors {
		err := c.pauseConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "error pausing connector %s", connectorName)
		}
	}

	return nil
}

func (c *ConnectorManager) pauseConnector(connectorName string) error {

	_, err := c.client.PauseConnector(connectorName)

	if err != nil {
		return err
	}

	return nil
}
