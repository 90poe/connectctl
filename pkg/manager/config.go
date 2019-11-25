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
}
