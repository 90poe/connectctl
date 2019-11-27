package manager

import (
	"net/http"
	"time"

	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/heptiolabs/healthcheck"
	"github.com/pkg/errors"
)

// ConnectorSource will return a slice of the desired connector configuration
type ConnectorSource func() ([]connect.Connector, error)

type client interface {
	CreateConnector(conn connect.Connector) (*http.Response, error)
	ListConnectors() ([]string, *http.Response, error)
	GetConnector(name string) (*connect.Connector, *http.Response, error)
	ListPlugins() ([]*connect.Plugin, *http.Response, error)
	GetConnectorStatus(name string) (*connect.ConnectorStatus, *http.Response, error)
	DeleteConnector(name string) (*http.Response, error)
	RestartConnectorTask(name string, taskID int) (*http.Response, error)
	UpdateConnectorConfig(name string, config connect.ConnectorConfig) (*connect.Connector, *http.Response, error)
	RestartConnector(name string) (*http.Response, error)
	ResumeConnector(name string) (*http.Response, error)
	PauseConnector(name string) (*http.Response, error)
}

// ConnectorManager manages connectors in a Kafka Connect cluster
type ConnectorManager struct {
	config *Config
	client client

	readinessState readinessState
}

// NewConnectorsManager creates a new ConnectorManager
func NewConnectorsManager(client client, config *Config) (*ConnectorManager, error) {
	//userAgent := fmt.Sprintf("90poe.io/connectctl/%s", config.Version)

	//client, err := connect.NewClient(config.ClusterURL, userAgent)
	//if err != nil {
	//	return nil, errors.Wrap(err, "error creating connect client")
	//}
	return &ConnectorManager{
		config:         config,
		client:         client,
		readinessState: unknownState,
	}, nil
}

type readinessState int

const (
	unknownState readinessState = iota
	okState
	errorState
)

// ReadinessCheck checks if we have been able to start syncing with kafka-connect
func (c *ConnectorManager) ReadinessCheck() (string, func() error) {
	return "connectctl-readiness-check", func() error {
		switch c.readinessState {
		case okState:
			return nil
		case unknownState, errorState:
			return errors.New("connectctl is not ready")
		}

		return nil
	}
}

// LivenessCheck checks if the the kafka-connect instance is running.
// The timeout of 2 seconds is arbitrary.
func (c *ConnectorManager) LivenessCheck() (string, func() error) {
	return "connectctl-liveness-check-kafka-connect-instance",
		healthcheck.HTTPGetCheck(c.config.ClusterURL, time.Second*2)
}
