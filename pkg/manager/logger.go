package manager

// Logger is the interface to implement to get all of the great news/updates
type Logger interface {
	Infof(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
}

type noopLogger struct{}

func (n *noopLogger) Infof(string, ...interface{}) {
	// do nothing
}
func (n *noopLogger) Warnf(string, ...interface{}) {
	// do nothing
}
func (n *noopLogger) Debugf(string, ...interface{}) {
	// do nothing
}
func (n *noopLogger) Errorf(string, ...interface{}) {
	// do nothing
}

func newNoopLogger() *noopLogger {
	return &noopLogger{}
}
