package connect

import "net/http"

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
