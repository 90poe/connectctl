package logging

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Configure sets up the logger
func Configure(logLevel string, logFile string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "getting hostname")
	}
	logrus.WithField("hostname", hostname)

	// use a file if you want
	if logFile != "" {
		f, errOpen := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND, 0660)
		if errOpen != nil {
			return errors.Wrapf(errOpen, "opening log file %s", logFile)
		}
		logrus.SetOutput(bufio.NewWriter(f))
	}

	if logLevel != "" {
		level, err := logrus.ParseLevel(strings.ToUpper(logLevel))
		if err != nil {
			return errors.Wrapf(err, "setting log level to %s", level)
		}
		logrus.SetLevel(level)
	}

	// always use the fulltimestamp
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return nil
}
