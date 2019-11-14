package manager

import (
	"time"
)

// Config represents the connect manager configuration
type Config struct {
	ClusterURL  string
	SyncPeriod  time.Duration
	AllowPurge  bool
	AutoRestart bool

	Version string
}
