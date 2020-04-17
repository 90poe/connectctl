package manager

import (
	"github.com/pkg/errors"

	"github.com/90poe/connectctl/pkg/client/connect"
)

// GetClusterInfo returns kafka cluster info
func (c *ConnectorManager) GetClusterInfo() (*connect.ClusterInfo, error) {
	clusterInfo, _, err := c.client.GetClusterInfo()

	if err != nil {
		return nil, errors.Wrap(err, "error getting cluster info")
	}

	return clusterInfo, nil
}
