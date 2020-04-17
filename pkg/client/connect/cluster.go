package connect

import (
	"net/http"
)

// ClusterInfo - this is new and not from the original author
type ClusterInfo struct {
	Version        string `json:"version"`
	Commit         string `json:"commit"`
	KafkaClusterID string `json:"kafka_cluster_id"`
}

// GetClusterInfo retrieves information about a cluster
//
// See: https://docs.confluent.io/current/connect/references/restapi.html#kconnect-cluster
func (c *Client) GetClusterInfo() (*ClusterInfo, *http.Response, error) {
	path := "/"
	clusterInfo := new(ClusterInfo)
	response, err := c.get(path, &clusterInfo)
	return clusterInfo, response, err
}
