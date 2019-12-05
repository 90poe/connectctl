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
	existing, _, err := c.client.ListConnectors()
	if err != nil {
		return errors.Wrap(err, "error listing connectors")
	}
	return c.resumeSpecifiedConnectors(existing)
}

func (c *ConnectorManager) resumeSpecifiedConnectors(connectors []string) error {
	for _, connectorName := range connectors {
		if _, err := c.client.ResumeConnector(connectorName); err != nil {
			return errors.Wrapf(err, "error resuming connector %s", connectorName)
		}
	}

	return nil
}
