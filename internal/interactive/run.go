package interactive

import (
	"context"

	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"

	"github.com/charmbracelet/huh"
)

type Options struct {
	LoadState func() (*state.State, error)
	Run       func() error
	Next      func() error
}

// Run starts the interactive loop driven by the kata state.
func Run(ctx context.Context, opts Options, translator i18n.Translator) error {
	for {
		st, err := opts.LoadState()
		if err != nil {
			return err
		}

		var choice Action

		selectField := huh.NewSelect[Action]().
			Title(translator.T("interactive.title_what_next")).
			Options(availableActions(st, translator)...).
			Value(&choice)

		note := huh.NewNote().Title(statusMessage(st, translator))
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
		}
	}
}
