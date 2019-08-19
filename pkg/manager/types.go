package manager

import (
	"github.com/90poe/connectctl/pkg/client/connect"
)

// ConnectorWithState is the connect config and state
type ConnectorWithState struct {
	Name           string                  `json:"name"`
	Config         connect.ConnectorConfig `json:"config,omitempty"`
	ConnectorState connect.ConnectorState  `json:"connectorState,omitempty"`
	Tasks          []connect.TaskState     `json:"tasks,omitempty"`
}
