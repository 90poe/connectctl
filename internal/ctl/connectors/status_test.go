package connectors

import (
	"testing"

	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/stretchr/testify/require"
)

func TestCountFailing(t *testing.T) {
	statusList := []*connect.ConnectorStatus{
		{
			Connector: connect.ConnectorState{
				State: "RUNNING",
			},
			Tasks: []connect.TaskState{
				{
					State: "RUNNING",
				},
				{
					State: "FAILED",
				},
			},
		},
		{
			Connector: connect.ConnectorState{
				State: "FAILED",
			},
			Tasks: []connect.TaskState{
				{
					State: "FAILED",
				},
				{
					State: "FAILED",
				},
			},
		},
	}

	connectorsFailing, tasksFailing := countFailing(statusList)

	require.Equal(t, 1, connectorsFailing)
	require.Equal(t, 3, tasksFailing)

}
