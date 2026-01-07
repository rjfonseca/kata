package interactive

import (
	"github.com/charmbracelet/huh"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
)

type Action int

const (
	ActionRun Action = iota
	ActionNext
	ActionExit
)

func availableActions(s *state.State, translator i18n.Translator) []huh.Option[Action] {
	var actions []huh.Option[Action]

	runOption := huh.Option[Action]{
		Key:   translator.T("interactive.action_run"),
		Value: ActionRun,
	}

	nextOption := huh.Option[Action]{
		Key:   translator.T("interactive.action_next"),
		Value: ActionNext,
	}

	exitOption := huh.Option[Action]{
		Key:   translator.T("interactive.action_exit"),
		Value: ActionExit,
	}

	switch {
	case s.TestPassing && !s.KataFinished:
		actions = append(actions, nextOption, runOption, exitOption)
	default:
		actions = append(actions, runOption, exitOption)
	}

	return actions
}
