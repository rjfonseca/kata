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
	},
	nil
}

// Run executes a task by name.
// It returns an error if the task fails or does not exist.
func (r *TaskRunner) Run(taskName string) error {
	ctx := context.Background()

	return r.executor.RunTask(ctx, &task.Call{
		Task: taskName,
	})
}

// TryRun executes a task by name if it exists.
// It returns an error only if the task exists but fails.
// If the task does not exist, it returns nil.
func (r *TaskRunner) TryRun(taskName string) error {
	err := r.Run(taskName)
	if err == nil {
		return nil
	}

	if isTaskNotFoundError(err) {
		return nil
	}

	return err
}

func isTaskNotFoundError(err error) bool {
	var notFoundErr *taskerrors.TaskNotFoundError
	return errors.As(err, &notFoundErr)
}

// ListTasks returns a list of all available tasks.
func (r *TaskRunner) ListTasks() ([]TaskInfo, error) {
	tasks, err := r.executor.GetTaskList()
	if err != nil {
		return nil, err
	}

	var result []TaskInfo
	for _, t := range tasks {
		// Filter out internal tasks if needed?
		// go-task has Internal bool field.
		if t.Internal {
			continue
		}
		result = append(result, TaskInfo{
			Name:        t.Task,
			Description: t.Desc,
		})
	}
	return result, nil
}

