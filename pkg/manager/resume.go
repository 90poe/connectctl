package manager

import "github.com/pkg/errors"

// Resume will resume a number of connectors in a cluster
func (c *ConnectorManager) Resume(connectors []string) error {
	if len(connectors) == 0 {
		return c.resumeAllConnectors()
	}

	return c.resumeSpecifiedConnectors(connectors)
}

func (c *ConnectorManager) resumeAllConnectors() error {
	c.logger.Info("resuming all connectors")

	existing, resp, err := c.client.ListConnectors()
	c.logger.WithField("response", resp).Trace("list connectors response")
	if err != nil {
		return errors.Wrap(err, "getting existing connectors")
	}

	for _, connectorName := range existing {
		err := c.resumeConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "resuming connector %s", connectorName)
		}
	}

	return nil
}

func (c *ConnectorManager) resumeSpecifiedConnectors(connectors []string) error {
	c.logger.Info("resuming specified connectors")
	for _, connectorName := range connectors {
		err := c.resumeConnector(connectorName)
		if err != nil {
			return errors.Wrapf(err, "resuming connector %s", connectorName)
		}
	}

	return nil
}

func (c *ConnectorManager) resumeConnector(connectorName string) error {
	connectLogger := c.logger.WithField("connector", connectorName)
	connectLogger.Info("resuming connector")

	resp, err := c.client.ResumeConnector(connectorName)
	connectLogger.WithField("response", resp).Trace("resume connector response")

	if err != nil {
		return errors.Wrap(err, "calling resume connector API")
	}

	connectLogger.Info("resumed connector")
	return nil
}
