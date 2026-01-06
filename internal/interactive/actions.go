package interactive

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rjfonseca/kata/internal/state"
)

type Action string

const (
	ActionRun  Action = "run_tests"
	ActionNext Action = "next_step"
	ActionExit Action = "exit"
)

const taskPrefix = "Task: "

func availableActions(s *state.State, customTasks []string) []huh.Option[Action] {
	var opts []huh.Option[Action]

	// Add Next Step if test passing and kata not finished
	if s.TestPassing && !s.KataFinished {
		opts = append(opts, ActionNext.HuhOption())
	}

	// Always add Run Tests
	opts = append(opts, ActionRun.HuhOption())

	// Add custom tasks
	for _, task := range customTasks {
		opts = append(opts, huh.Option[Action]{
			Key:   fmt.Sprintf("%s%s", taskPrefix, task),
			Value: Action("task:" + task),
		})
	}

	// Always add Exit
	opts = append(opts, ActionExit.HuhOption())

	return opts
}

func (a Action) String() string {
	switch a {
	case ActionRun:
		return "Run tests"
	case ActionNext:
		return "Next step"
	case ActionExit:
		return "Exit"
	default:
		if strings.HasPrefix(string(a), "task:") {
			return fmt.Sprintf("%s%s", taskPrefix, strings.TrimPrefix(string(a), "task:"))
		}
		return "Unknown"
	}
}

func (a Action) HuhOption() huh.Option[Action] {
	return huh.Option[Action]{
		Key:   a.String(),
		Value: a,
	}
}
