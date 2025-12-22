package taskrunner

// TaskInfo describes a available task.
type TaskInfo struct {
	Name        string
	Description string
}

// Runner executes named tasks defined by the project.
type Runner interface {
	Run(taskName string) error
	TryRun(taskName string) error
	ListTasks() ([]TaskInfo, error)
}
