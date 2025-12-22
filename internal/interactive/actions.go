package interactive

import (
	"github.com/charmbracelet/huh"
	"github.com/rjfonseca/kata/internal/state"
)

type Action int

const (
	ActionRun Action = iota
	ActionNext
	ActionExit
)

func availableActions(s *state.State) []huh.Option[Action] {
	switch {
	case s.TestPassing && !s.KataFinished:
		return []huh.Option[Action]{
			ActionNext.HuhOption(),
			ActionRun.HuhOption(),
			ActionExit.HuhOption(),
		}
	default:
		return []huh.Option[Action]{
			ActionRun.HuhOption(),
			ActionExit.HuhOption(),
		}
	}
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
		return "Unknown"
	}
}

func (a Action) HuhOption() huh.Option[Action] {
	return huh.Option[Action]{
		Key:   a.String(),
		Value: a,
	}
}
