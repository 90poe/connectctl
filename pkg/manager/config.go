package manager

import (
	"time"
)

// Config represents the connect manager configuration
type Config struct {
	ClusterURL  string        `json:"cluster_url"`
	SyncPeriod  time.Duration `json:"sync_period"`
	AllowPurge  bool          `json:"allow_purge"`
	AutoRestart bool          `json:"auto_restart"`

	Version string `json:"version"`

	RestartPolicy *RestartPolicy `json:"restart_policy"`
}

// RestartPolicy lists each connectors maximum restart policy
// If a policy does not exist for a connector the connector ir task will be restarted once.
// If a connector or task is restarted the count of failed attempts is reset.
// If the maximum number of unsuccessful restarts is reached the manager will
// return and connectctl will stop.
type RestartPolicy struct {
	Connectors map[string]Policy `json:"connectors"`
}

// Policy contains a collection of values to be managed
type Policy struct {
	MaxConnectorRestarts   int           `json:"max_connector_restarts"`
	ConnectorRestartPeriod time.Duration `json:"connector_restart_period"`
	MaxTaskRestarts        int           `json:"max_task_restarts"`
	TaskRestartPeriod      time.Duration `json:"task_restart_period"`
}
