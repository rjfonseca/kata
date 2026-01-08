package taskrunner

// Runner executes named tasks defined by the project.
type Runner interface {
	Run(taskName string) error
	RunOptional(taskName string) error
}
