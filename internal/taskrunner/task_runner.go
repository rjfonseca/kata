package taskrunner

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-task/task/v3"
	taskerrors "github.com/go-task/task/v3/errors"
)

// TaskRunner executes tasks using go-task as a library.
type TaskRunner struct {
	executor *task.Executor
}

// New creates a TaskRunner bound to a working directory.
func New(root string) (*TaskRunner, error) {
	e := task.NewExecutor(task.WithVersionCheck(true))
	e.Dir = root // Set the directory for the executor
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

// RunOptional executes a task by name if it exists.
// It returns nil if the task does not exist.
func (r *TaskRunner) RunOptional(taskName string) error {
	err := r.Run(taskName)
	if err == nil {
		return nil
	}

	var notFoundErr *taskerrors.TaskNotFoundError
	if errors.As(err, &notFoundErr) {
		return nil
	}

	return err
}
