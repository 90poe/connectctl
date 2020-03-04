package manager

import (
	"net/http"
	"testing"

	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_Manage_MissingConnectorsAreAdded(t *testing.T) {

	createCalled := false

	mock := &mocks.FakeClient{
		GetConnectorStub: func(string) (*connect.Connector, *http.Response, error) {
			return nil, nil, &connect.APIError{Code: http.StatusNotFound}
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

func Test_Manage_ExistingConnectorsAreRemovedIfNotListed(t *testing.T) {

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

func Test_Manage_ErrorsAreAPIErrorsIfUnwrapped(t *testing.T) {

	mock := &mocks.FakeClient{
		GetConnectorStub: func(string) (*connect.Connector, *http.Response, error) {
			return nil, nil, &connect.APIError{Code: http.StatusInternalServerError}
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
	rootCause := errors.Cause(err)
	require.True(t, connect.IsAPIError(rootCause))
}

func Test_Manage_ConnectorRunning_FailedTasksAreRestarted(t *testing.T) {

	mock := &mocks.FakeClient{
		GetConnectorStatusStub: func(string) (*connect.ConnectorStatus, *http.Response, error) {
			return &connect.ConnectorStatus{
				Connector: connect.ConnectorState{
					State: "RUNNING",
				},
				Tasks: []connect.TaskState{
					connect.TaskState{
						State: "FAILED",
					},
				},
			}, nil, nil
		}}

	config := &Config{
		AutoRestart: true,
	}

	cm, err := NewConnectorsManager(mock, config)
	require.Nil(t, err)

	source := func() ([]connect.Connector, error) {
		return []connect.Connector{
			connect.Connector{Name: "foo"},
		}, nil
	}

	err = cm.Sync(source)
	require.Nil(t, err)
	require.Equal(t, mock.RestartConnectorTaskCallCount(), 1)
}

func Test_Manage_ConnectorFailed_IsRestarted(t *testing.T) {

	mock := &mocks.FakeClient{
		GetConnectorStatusStub: func(string) (*connect.ConnectorStatus, *http.Response, error) {
			return &connect.ConnectorStatus{
				Connector: connect.ConnectorState{
					State: "FAILED",
				},
				Tasks: []connect.TaskState{
					connect.TaskState{
						State: "FAILED",
					},
				},
			}, nil, nil
		},
	}

	config := &Config{
		AutoRestart: true,
	}

	cm, err := NewConnectorsManager(mock, config)
	require.Nil(t, err)

	source := func() ([]connect.Connector, error) {
		return []connect.Connector{
			connect.Connector{Name: "foo"},
		}, nil
	}

	err = cm.Sync(source)
	require.Nil(t, err)
	require.Equal(t, mock.RestartConnectorTaskCallCount(), 0)
}
