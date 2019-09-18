package manager

import (
	"fmt"

	"github.com/90poe/connectctl/pkg/client/connect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ConnectorSource will return a slice of the desired connector configuration
type ConnectorSource func() ([]*connect.Connector, error)

// ConnectorManager manages connectors in a Kafka Connect cluster
type ConnectorManager struct {
	config *Config
	client *connect.Client
	logger *log.Entry
}

// NewConnectorsManager creates a new ConnectorManager
func NewConnectorsManager(config *Config) (*ConnectorManager, error) {
	userAgent := fmt.Sprintf("90poe.io/connectctl/%s", config.Version)

	client, err := connect.NewClient(config.ClusterURL, userAgent)
	if err != nil {
		return nil, errors.Wrap(err, "creating connect client")
	}

	return &ConnectorManager{
		config: config,
		client: client,
		logger: log.WithField("cluster", config.ClusterURL),
	}, nil
}
