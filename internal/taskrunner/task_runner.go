package taskrunner

import (
	"context"
	"fmt"

	"github.com/go-task/task/v3"
)

// TaskRunner executes tasks using go-task as a library.
type TaskRunner struct {
	executor *task.Executor
}

// New creates a TaskRunner bound to a working directory.
func New() (*TaskRunner, error) {
	e := task.NewExecutor(task.WithVersionCheck(true))
	err := e.Setup()
	if err != nil {
		return nil, fmt.Errorf("creating Taskfile executor: %w", err)
	}
	return &TaskRunner{
		executor: e,
	}, nil

}

// Run executes a task by name.
// It returns an error if the task fails or does not exist.
func (r *TaskRunner) Run(taskName string) error {
	ctx := context.Background()

	return r.executor.RunTask(ctx, &task.Call{
		Task: taskName,
	})
}
