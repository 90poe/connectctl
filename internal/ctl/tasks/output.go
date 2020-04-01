package tasks

import (
	"bytes"
	"encoding/json"
	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/jedib0t/go-pretty/table"
	"github.com/pkg/errors"
)

type OutputType string

const (
	OutputTypeJSON  OutputType = "json"
	OutputTypeTable OutputType = "table"
)

type TaskListOutputFn func([]connect.Task) ([]byte, error)

func OutputTaskListAsJSON(tasks []connect.Task) ([]byte, error) {
	b, err := DefaultMarshalIndent(tasks)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal tasks")
	}
	return b, nil
}

func OutputTaskListAsTable(tasks []connect.Task) ([]byte, error) {
	var buff bytes.Buffer
	t := table.NewWriter()
	t.Style().Options.SeparateRows = true
	t.SetOutputMirror(&buff)
	t.AppendHeader(table.Row{"Connector Name", "ID", "Config"})
	for _, task := range tasks {
		configBytes, err := DefaultMarshalIndent(task.Config)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal config %v", task.Config)
		}
		t.AppendRow(table.Row{task.ID.ConnectorName, task.ID.ID, string(configBytes)})
	}
	t.Render()
	return buff.Bytes(), nil
}

func DefaultMarshalIndent(value interface{}) ([]byte, error) {
	return json.MarshalIndent(value, "", "	")
}

type TaskStateOutputFn func(*connect.TaskState) ([]byte, error)

func OutputTaskStateAsJSON(taskState *connect.TaskState) ([]byte, error) {
	b, err := DefaultMarshalIndent(taskState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal task state")
	}
	return b, nil
}

func OutputTaskStateAsTable(taskState *connect.TaskState) ([]byte, error) {
	var buff bytes.Buffer
	t := table.NewWriter()
	t.SetOutputMirror(&buff)
	t.AppendHeader(table.Row{"ID", "WorkerID", "State", "Trace"})
	t.AppendRow(table.Row{taskState.ID, taskState.WorkerID, taskState.State, taskState.Trace})
	t.Render()
	return buff.Bytes(), nil
}
