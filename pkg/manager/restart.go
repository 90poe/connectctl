package manager

import "github.com/pkg/errors"

func (c *ConnectorManager) Restart(connectors []string) error {
	if len(connectors) == 0 {
		return c.restartAllConnectors()
	}

	return c.restartSpecifiedConnectors(connectors)
}

func (c *ConnectorManager) restartAllConnectors() error {
	c.logger.Info("restarting all connectors")

	existing, resp, err := c.client.ListConnectors()
	c.logger.WithField("response", resp).Trace("list connectors response")
	if err != nil {
		return errors.Wrap(err, "getting existing connectors")
	}

	for _, connectorName := range existing {
		err := c.restartConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "restarting connector %s", connectorName)
		}
	}

	return nil
}

func (c *ConnectorManager) restartSpecifiedConnectors(connectors []string) error {
	c.logger.Info("restarting specified connectors")
	for _, connectorName := range connectors {
		err := c.restartConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "restarting connector %s", connectorName)
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
