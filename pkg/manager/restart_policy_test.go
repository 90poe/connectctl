package manager

import (
	"testing"
	"time"

	"github.com/90poe/connectctl/pkg/client/connect"

	"github.com/stretchr/testify/require"
)

func Test_RestartPolicy_Default(t *testing.T) {
	t.Parallel()

	connectors := []connect.Connector{
		connect.Connector{Name: "foo"},
	}

	policy := runtimePolicyFromConnectors(connectors, nil)

	require.Len(t, policy, 1)
	require.NotNil(t, policy["foo"])

	foo := policy["foo"]

	require.Equal(t, 1, foo.MaxConnectorRestarts)
	require.Equal(t, 1, foo.MaxTaskRestarts)
	require.Equal(t, defaultRestartPeriod, foo.ConnectorRestartPeriod)
	require.Equal(t, defaultRestartPeriod, foo.TaskRestartPeriod)
}

func Test_RestartPolicy_Globals(t *testing.T) {
	t.Parallel()

	connectors := []connect.Connector{
		connect.Connector{Name: "foo"},
	}

	policy := runtimePolicyFromConnectors(connectors, &Config{
		GlobalMaxConnectorRestarts:   97,
		GlobalConnectorRestartPeriod: time.Second * 98,
		GlobalMaxTaskRestarts:        99,
		GlobalTaskRestartPeriod:      time.Second * 100,
	})

	require.Len(t, policy, 1)
	require.NotNil(t, policy["foo"])

	foo := policy["foo"]

	require.Equal(t, 97, foo.MaxConnectorRestarts)
	require.Equal(t, time.Second*98, foo.ConnectorRestartPeriod)
	require.Equal(t, 99, foo.MaxTaskRestarts)
	require.Equal(t, time.Second*100, foo.TaskRestartPeriod)
}

func Test_RestartPolicy_Override(t *testing.T) {
	t.Parallel()

	connectors := []connect.Connector{
		connect.Connector{Name: "foo"},
	}

	ovveride := RestartPolicy{
		Connectors: map[string]Policy{
			"foo": Policy{
				MaxConnectorRestarts:   10,
				MaxTaskRestarts:        11,
				TaskRestartPeriod:      time.Second * 100,
				ConnectorRestartPeriod: time.Second * 101,
			},
		},
	}

	config := &Config{RestartOverrides: &ovveride}
	policy := runtimePolicyFromConnectors(connectors, config)

	require.Len(t, policy, 1)
	require.NotNil(t, policy["foo"])

	foo := policy["foo"]

	require.Equal(t, 10, foo.MaxConnectorRestarts)
	require.Equal(t, 11, foo.MaxTaskRestarts)
	require.Equal(t, time.Second*101, foo.ConnectorRestartPeriod)
	require.Equal(t, time.Second*100, foo.TaskRestartPeriod)
}
