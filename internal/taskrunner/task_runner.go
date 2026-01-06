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

// ListTasks returns a list of all available tasks.
func (r *TaskRunner) ListTasks() ([]string, error) {
	var tasks []string
	// Keys(nil) returns an iterator of task names in sorted order.
	for task := range r.executor.Taskfile.Tasks.Keys(nil) {
		tasks = append(tasks, task)
	}
	return tasks, nil
}
