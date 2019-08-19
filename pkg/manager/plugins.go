package manager

import (
	"github.com/pkg/errors"

	"github.com/90poe/connectctl/pkg/client/connect"
)

// GetAllPlugins returns all the connector plugins installed
func (c *ConnectorManager) GetAllPlugins() ([]*connect.Plugin, error) {
	c.logger.Debug("getting all connector plugins")

	plugins, resp, err := c.client.ListPlugins()
	c.logger.WithField("response", resp).Trace("list plugins response")
	if err != nil {
		return nil, errors.Wrap(err, "list plugins with api")
	}

	return plugins, nil
}
