package manager

import (
	"time"
)

// Config represents the connect manager configuration
type Config struct {
	ClusterURL        string        `json:"cluster_url"`
	SyncPeriod        time.Duration `json:"sync_period"`
	InitialWaitPeriod time.Duration `json:"initial_wait_period"`
	AllowPurge        bool          `json:"allow_purge"`
	AutoRestart       bool          `json:"auto_restart"`

	Version string `json:"version"`

	GlobalConnectorRestartsMax   int           `json:"global_connector_restarts_max"`
	GlobalConnectorRestartPeriod time.Duration `json:"global_connector_restart_period"`
	GlobalTaskRestartsMax        int           `json:"global_task_restarts_max"`
	GlobalTaskRestartPeriod      time.Duration `json:"global_task_restart_period"`

	RestartOverrides *RestartPolicy `json:"restart_policy"`
}

// RestartPolicy lists each connectors maximum restart policy
// If AutoRestart == true
// If a policy does not exist for a connector the connector or task will be restarted once.
// If a connector or task is restarted the count of failed attempts is reset.
// If the number of unsuccessful restarts is reached the manager will return and connectctl will stop.
type RestartPolicy struct {
	Connectors map[string]Policy `json:"connectors"`
}

// Policy contains a collection of values to be managed
type Policy struct {
	ConnectorRestartsMax   int           `json:"connector_restarts_max"`
	ConnectorRestartPeriod time.Duration `json:"connector_restart_period"`
	TaskRestartsMax        int           `json:"task_restarts_max"`
	TaskRestartPeriod      time.Duration `json:"task_restart_period"`
}
