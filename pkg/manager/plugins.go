package manager

import (
	"github.com/pkg/errors"

	"github.com/90poe/connectctl/pkg/client/connect"
)

// GetAllPlugins returns all the connector plugins installed
func (c *ConnectorManager) GetAllPlugins() ([]*connect.Plugin, error) {
	plugins, _, err := c.client.ListPlugins()

	if err != nil {
		return nil, errors.Wrap(err, "error listing plugins")
	}

	return plugins, nil
}

// ValidatePlugins returns validation results of a connector config
func (c *ConnectorManager) ValidatePlugins(config connect.ConnectorConfig) (*connect.ConfigValidation, error) {
	validation, _, err := c.client.ValidatePlugins(config)

	if err != nil {
		return nil, errors.Wrap(err, "error validating plugins")
	}

	return validation, nil
}
