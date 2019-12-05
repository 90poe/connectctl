package manager

import (
	"testing"

	"github.com/90poe/connectctl/pkg/client/connect"
	"github.com/stretchr/testify/assert"
)

func TestTasks_ByID_Found(t *testing.T) {
	f := ByID(3, 2, 1)
	assert.True(t, f(connect.TaskState{ID: 1}))
	assert.True(t, f(connect.TaskState{ID: 2}))
	assert.True(t, f(connect.TaskState{ID: 3}))
	assert.False(t, f(connect.TaskState{ID: 0}))
	assert.False(t, f(connect.TaskState{ID: 4}))
}

func TestTasks_IsRunning(t *testing.T) {
	f := IsRunning
	assert.True(t, f(connect.TaskState{State: "RUNNING"}))
	assert.False(t, f(connect.TaskState{State: "OTHER"}))
}

func TestTasks_IsNotRunning(t *testing.T) {
	f := IsNotRunning
	assert.False(t, f(connect.TaskState{State: "RUNNING"}))
	assert.True(t, f(connect.TaskState{State: "OTHER"}))
}

func TestTasks_IDsEmpty(t *testing.T) {
	var tasks Tasks
	assert.Equal(t, []int{}, tasks.IDs())
}

func TestTasks_IDsNotEmpty(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3},
		Tasks{
			connect.TaskState{ID: 1},
			connect.TaskState{ID: 2},
			connect.TaskState{ID: 3},
		}.IDs())
}

func TestTasks_Filter(t *testing.T) {
	assert.Equal(t, Tasks{
		connect.TaskState{ID: 2},
		connect.TaskState{ID: 1},
		connect.TaskState{ID: 0},
	},
		Tasks{
			connect.TaskState{ID: -2},
			connect.TaskState{ID: 2},
			connect.TaskState{ID: -1},
			connect.TaskState{ID: 1},
			connect.TaskState{ID: 0},
		}.Filter(func(task connect.TaskState) bool {
			return task.ID >= 0
		}))
}
