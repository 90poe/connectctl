package connect

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// This is new and not from the original author

// Plugin represents a Kafka Connect connector plugin
type Plugin struct {
	Class   string `json:"class"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

type FieldDefinition struct {
	Name          string                   `json:"name"`
	Type          string                   `json:"type"`
	Required      bool                     `json:"required"`
	DefaultValue  *string                  `json:"default_value"`
	Importance    string                   `json:"importance"`
	Documentation string                   `json:"documentation"`
	Group         string                   `json:"group"`
	Width         string                   `json:"width"`
	DisplayName   string                   `json:"display_name"`
	Dependents    []map[string]interface{} `json:"dependents"` //unknown type
	Order         int                      `json:"order"`
}

type FieldValue struct {
	Name              string    `json:"name"`
	Value             *string   `json:"value"`
	RecommendedValues []*string `json:"recommended_values"`
	Errors            []string  `json:"errors"`
	Visible           bool      `json:"visible"`
}

type FieldValidation struct {
	Definition FieldDefinition `json:"definition"`
	Value      FieldValue      `json:"value"`
}

type ConfigValidation struct {
	Name       string            `json:"name"`
	ErrorCount int               `json:"error_count"`
	Groups     []string          `json:"groups"`
	Configs    []FieldValidation `json:"configs"`
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
func (c *Client) ValidatePlugins(config ConnectorConfig) (*ConfigValidation, *http.Response, error) {
	connectorClass, ok := config["connector.class"]
	if !ok {
		return nil, nil, errors.New("missing required key in config: 'connector.class'")
	}

	tuple := strings.Split(connectorClass, ".")
	path := fmt.Sprintf("connector-plugins/%s/config/validate", tuple[len(tuple)-1])

	var validation ConfigValidation
	response, err := c.doRequest("PUT", path, config, &validation)

	return &validation, response, err
}
