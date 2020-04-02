package tasks

import (
	"github.com/90poe/connectctl/pkg/client/connect"
)

func findTaskByID(tasks []connect.Task, id int) (connect.Task, bool) {
	for _, t := range tasks {
		if t.ID.ID == id {
			return t, true
		}
	}
	return connect.Task{}, false
}
