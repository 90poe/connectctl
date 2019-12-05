package manager

import (
	"net/http"
	"testing"

	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager/mocks"
	"github.com/stretchr/testify/require"
)

func Test_MissingConnectorsAreAdded(t *testing.T) {

	createCalled := false

	mock := &mocks.FakeClient{
		GetConnectorStub: func(string) (*connect.Connector, *http.Response, error) {
			return nil, nil, connect.APIError{Code: http.StatusNotFound}
		},
		CreateConnectorStub: func(c connect.Connector) (*http.Response, error) {

			require.Equal(t, c.Name, "one")

			createCalled = true
			return nil, nil
		},
	}
	config := &Config{}

	cm, err := NewConnectorsManager(mock, config)
	require.Nil(t, err)

	source := func() ([]connect.Connector, error) {
		return []connect.Connector{
			connect.Connector{Name: "one"},
		}, nil

	}

	err = cm.Sync(source)
	require.Nil(t, err)
	require.True(t, createCalled)

}

func Test_ExistingConnectorsAreRemovedIfNotListed(t *testing.T) {

	deleteCalled := false

	mock := &mocks.FakeClient{
		ListConnectorsStub: func() ([]string, *http.Response, error) {
			return []string{"delete-me"}, nil, nil
		},
		DeleteConnectorStub: func(string) (*http.Response, error) {
			deleteCalled = true
			return nil, nil
		},
	}
	config := &Config{
		AllowPurge: true,
	}

	cm, err := NewConnectorsManager(mock, config)
	require.Nil(t, err)

	source := func() ([]connect.Connector, error) {
		return []connect.Connector{}, nil
	}

	err = cm.Sync(source)
	require.Nil(t, err)
	require.True(t, deleteCalled)

}
