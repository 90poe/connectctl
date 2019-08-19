// +build !windows

package signals

import (
	"os"
	"syscall"
)

//nolint:gochecknoglobals
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
