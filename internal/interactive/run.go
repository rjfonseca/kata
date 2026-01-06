package interactive

import (
	"context"
	"strings"

	"github.com/rjfonseca/kata/internal/state"

	"github.com/charmbracelet/huh"
)

type Options struct {
	LoadState func() (*state.State, error)

	Run         func() error
	Next        func() error
	RunTask     func(string) error
	CustomTasks []string
}

// Run starts the interactive loop driven by the kata state.
func Run(ctx context.Context, opts Options) error {
	for {
		st, err := opts.LoadState()
		if err != nil {
			return err
		}

		var choice Action

		selectField := huh.NewSelect[Action]().
			Title("What would you like to do next?").
			Options(availableActions(st, opts.CustomTasks)...).
			Value(&choice)

		note := huh.NewNote().Title(statusMessage(st))
		form := huh.NewForm(
			huh.NewGroup(note, selectField),
		)

		if err := form.Run(); err != nil {
			// Ctrl+C or Esc exits interactive mode cleanly
			return nil
		}

		switch choice {
		case ActionRun:
			_ = opts.Run()
		case ActionNext:
			_ = opts.Next()
		case ActionExit:
			return nil
		default:
			if strings.HasPrefix(string(choice), "task:") {
				taskName := strings.TrimPrefix(string(choice), "task:")
				if opts.RunTask != nil {
					_ = opts.RunTask(taskName)
				}
			}
		}
	}
}
