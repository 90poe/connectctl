package manager

import (
	"net/http"
	"testing"
	"time"

	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/90poe/connectctl/pkg/manager/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_Manage_MissingConnectorsAreAdded(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	require.NotNil(t, err)
	require.Equal(t, 2, mock.RestartConnectorTaskCallCount())
}

func Test_Manage_ConnectorFailed_IsRestarted(t *testing.T) {
	t.Parallel()
	count := 0

	mock := &mocks.FakeClient{
		GetConnectorStatusStub: func(string) (*connect.ConnectorStatus, *http.Response, error) {
			if count == 0 {
				count++
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
			} else {
				return &connect.ConnectorStatus{
					Connector: connect.ConnectorState{
						State: "RUNNING",
					},
					Tasks: []connect.TaskState{
						connect.TaskState{
							State: "RUNNING",
						},
					},
				}, nil, nil
			}
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
	require.Equal(t, 1, mock.RestartConnectorCallCount())
}

func Test_Manage_ConnectorFailed_IsRestarted_WithPolicy(t *testing.T) {
	t.Parallel()

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
		RestartConnectorStub: func(string) (*http.Response, error) {
			return nil, nil
		},
	}

	config := &Config{
		AutoRestart: true,
		RestartOverrides: &RestartPolicy{
			Connectors: map[string]Policy{
				"foo": Policy{
					ConnectorRestartsMax:   10,
					ConnectorRestartPeriod: time.Millisecond,
				},
			},
		},
	}

	cm, err := NewConnectorsManager(mock, config)
	require.Nil(t, err)

	source := func() ([]connect.Connector, error) {
		return []connect.Connector{
			connect.Connector{Name: "foo"},
		}, nil
	}

	err = cm.Sync(source)
	require.NotNil(t, err)
	require.Equal(t, 11, mock.RestartConnectorCallCount())
	require.Equal(t, 11, mock.GetConnectorStatusCallCount())
}

func Test_Manage_ConnectorFailed_IsRestarted_WithPolicy_RestartWorks(t *testing.T) {
	t.Parallel()
	count := 0

	mock := &mocks.FakeClient{
		GetConnectorStatusStub: func(string) (*connect.ConnectorStatus, *http.Response, error) {
			if count == 0 {
				count++
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
			} else {
				return &connect.ConnectorStatus{
					Connector: connect.ConnectorState{
						State: "RUNNING",
					},
					Tasks: []connect.TaskState{
						connect.TaskState{
							State: "RUNNING",
						},
					},
				}, nil, nil

			}
		},
		RestartConnectorStub: func(string) (*http.Response, error) {
			return nil, nil
		},
	}

	config := &Config{
		AutoRestart: true,
		RestartOverrides: &RestartPolicy{
			Connectors: map[string]Policy{
				"foo": Policy{
					ConnectorRestartsMax:   10,
					ConnectorRestartPeriod: time.Millisecond,
					TaskRestartsMax:        0,
					TaskRestartPeriod:      time.Millisecond,
				},
			},
		},
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
	require.Equal(t, 1, mock.RestartConnectorCallCount())
	require.Equal(t, 3, mock.GetConnectorStatusCallCount())
}
