package tasks

import (
	"github.com/90poe/connectctl/pkg/client/connect"
	"reflect"
	"testing"
)

func TestFindTaskByID(t *testing.T) {
	f := func(tasks []connect.Task, id int, expected connect.Task, found bool) {
		t.Helper()
		task, ok := findTaskByID(tasks, id)
		if ok != found {
			t.Fatalf("expected '%t', got '%t'", found, ok)
		}
		if !reflect.DeepEqual(task, expected) {
			t.Fatalf("expected %#v, got %#v", expected, task)
		}
	}

	empty := connect.Task{}

	tasks := []connect.Task{
		{
			ID: connect.TaskID{
				ConnectorName: "a",
				ID:            1,
			},
			Config: nil,
		},
		{
			ID: connect.TaskID{
				ConnectorName: "b",
				ID:            3,
			},
			Config: nil,
		},
		{
			ID: connect.TaskID{
				ConnectorName: "c",
				ID:            2,
			},
			Config: nil,
		},
	}

	f(nil, 1, empty, false)
	f([]connect.Task{}, 1, empty, false)
	f(
		tasks,
		1,
		connect.Task{
			ID: connect.TaskID{
				ConnectorName: "a",
				ID:            1,
			},
			Config: nil,
		},
		true,
	)
	f(
		tasks,
		3,
		connect.Task{
			ID: connect.TaskID{
				ConnectorName: "b",
				ID:            3,
			},
			Config: nil,
		},
		true,
	)
	f(
		tasks,
		2,
		connect.Task{
			ID: connect.TaskID{
				ConnectorName: "c",
				ID:            2,
			},
			Config: nil,
		},
		true,
	)
	f(tasks, -1, empty, false)
	f(tasks, 4, empty, false)
}
