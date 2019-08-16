package manager

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Config represent the connect manager configuration
type Config struct {
	ClusterURL string
	SyncPeriod time.Duration
	AllowPurge bool

	Logger *log.Entry
}
