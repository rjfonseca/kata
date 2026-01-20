package interactive

import (
	"github.com/charmbracelet/huh"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

type Action string

const (
	ActionRun  Action = "run"
	ActionNext Action = "next"
	ActionExit Action = "exit"
)

func ActionTask(name string) Action {
	return Action("task:" + name)
}

func availableActions(s *state.State, tasks []taskrunner.TaskInfo) []huh.Option[Action] {
	var opts []huh.Option[Action]

	if s.TestPassing && !s.KataFinished {
		opts = append(opts, huh.NewOption("Next step", ActionNext))
	}

	opts = append(opts, huh.NewOption("Run tests", ActionRun))

	for _, t := range tasks {
		opts = append(opts, huh.NewOption("Task: "+t.Name, ActionTask(t.Name)))
	}

	opts = append(opts, huh.NewOption("Exit", ActionExit))
	return opts
}
