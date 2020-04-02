package tasks

import (
	"github.com/pkg/errors"
	"testing"
)

func TestGenericOptions_Validate(t *testing.T) {
	f := func(opts GenericOptions, expectedErr error) {
		err := opts.Validate()
		if err != nil && errors.Cause(err) == expectedErr {
			t.Fatalf("expected %#v, got %#v", expectedErr, err)
		}
	}

	const (
		clusterURL    = "http://localhost:8083"
		connectorName = "connector"
	)

	var (
		clusterErr   = errors.New("--cluster is not a valid URI")
		connectorErr = errors.New("--connector name is empty")
	)

	f(
		GenericOptions{
			ClusterURL: "",
		},
		clusterErr,
	)
	f(
		GenericOptions{
			ConnectorName: "",
		},
		clusterErr,
	)
	f(
		GenericOptions{
			ClusterURL: "simple:string",
		},
		clusterErr,
	)
	f(
		GenericOptions{
			ClusterURL: clusterURL,
		},
		connectorErr,
	)
	f(
		GenericOptions{
			ClusterURL:    clusterURL,
			ConnectorName: connectorName,
		},
		nil,
	)
	f(
		GenericOptions{
			ClusterURL:    "www.google.com",
			ConnectorName: connectorName,
		},
		clusterErr,
	)
	f(
		GenericOptions{
			ClusterURL:    "https://www.google.com",
			ConnectorName: connectorName,
		},
		nil,
	)
}
