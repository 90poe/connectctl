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
	c.logger.Info("pausing all connectors")

	existing, resp, err := c.client.ListConnectors()
	c.logger.WithField("response", resp).Trace("list connectors response")
	if err != nil {
		return errors.Wrap(err, "getting existing connectors")
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
	c.logger.Info("pausing specified connectors")
	for _, connectorName := range connectors {
		err := c.pauseConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "pausing connector %s", connectorName)
		}
	}

	return nil
}

func (c *ConnectorManager) pauseConnector(connectorName string) error {
	connectLogger := c.logger.WithField("connector", connectorName)
	connectLogger.Info("pausing connector")

	resp, err := c.client.PauseConnector(connectorName)
	connectLogger.WithField("response", resp).Trace("pause connector response")

	if err != nil {
		return errors.Wrap(err, "calling pause connector API")
	}

	connectLogger.Info("paused connector")
	return nil
}
