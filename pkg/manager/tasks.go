package manager

import (
	"github.com/90poe/connectctl/pkg/client/connect"
	"sort"
)

type Tasks []connect.TaskState

// TaskPredicate is a function that performs some test on a connect.TaskState
type TaskPredicate func(*connect.TaskState) bool

// IsRunning returns true if the connector task is in a RUNNING state
func IsRunning(task *connect.TaskState) bool {
	return task.State == "RUNNING"
}

// IsNotRunning returns true if the connector task is not in a RUNNING state
func IsNotRunning(task *connect.TaskState) bool {
	return task.State != "RUNNING"
}

// ByID returns a predicate that returns true if the connector task has one of the given task IDs
func ByID(taskIDs ...int) TaskPredicate {
	sort.Ints(taskIDs)
	return func(task *connect.TaskState) bool {
		found := sort.SearchInts(taskIDs, task.ID)
		return found < len(taskIDs) && taskIDs[found] == task.ID
	}
}

// IDs returns a subset of the Tasks for which the predicate returns true
func (t Tasks) Filter(predicate TaskPredicate) Tasks {
	var found Tasks
	for _, task := range t {
		if predicate(&task) {
			found = append(found, task)
		}
	}
	return found
}

// IDs returns the slice of task IDs
func (t Tasks) IDs() []int {
	taskIDs := make([]int, len(t))
	for i, task := range t {
		taskIDs[i] = task.ID
	}
	return taskIDs
}
