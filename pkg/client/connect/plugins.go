package connect

import (
	"errors"
	"net/http"
)

// This is new and not from the original author

// Plugin represents a Kafka Connect connector plugin
type Plugin struct {
	Class   string `json:"class"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

// ListPlugins retrieves a list of the installed plugins.
// Note that the API only checks for connectors on the worker
// that handles the request, which means it is possible to see
// inconsistent results, especially during a rolling upgrade if
// you add new connector jars
// See: https://docs.confluent.io/current/connect/references/restapi.html#get--connector-plugins-
func (c *Client) ListPlugins() ([]*Plugin, *http.Response, error) {
	path := "connector-plugins"
	var names []*Plugin

	response, err := c.get(path, &names)
	return names, response, err
}

// ValidatePlugins validates the provided configuration values against the configuration definition.
// See: https://docs.confluent.io/current/connect/references/restapi.html#put--connector-plugins-(string-name)-config-validate
func (c *Client) ValidatePlugins(config ConnectorConfig) ([]*Plugin, *http.Response, error) {
	connectorClass, ok := config["connector.class"]
	if !ok {
		return nil, nil, errors.New("missing required key in config: connector.class")
	}
	// TODO
	path := "connector-plugins"
	var names []*Plugin

	response, err := c.get(path, &names)
	return names, response, err
}
