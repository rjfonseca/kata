package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/scaffold"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

func Next(root string, stateRepo *state.Repository, runner taskrunner.Runner, translator i18n.Translator) error {
	s, err := stateRepo.Load()
	if err != nil {
		return err
	}

	// Case 1: tests are failing (Red)
	if !s.TestPassing {
		return errors.New(translator.T("next.error_tests_failing"))
	}

	// Case 2: there is a next step to advance to
	if s.HasNextStep() {
		if err := s.Next(); err != nil {
			return err
		}

		stepName := s.Steps[s.CurrentStepIndex]
		slog.Info(translator.T("next.log_advancing"), "step", stepName)

		manifest, err := scaffold.LoadManifest(root)
		if err != nil {
			return err
		}

		copier := &scaffold.Copier{
			Root:     root,
			Manifest: manifest,
		}

		stepDir := filepath.Join(
			".kata",
			"catalog",
			s.KataName,
			"steps",
			stepName,
		)

		if err := copier.Apply(os.DirFS(stepDir), "step:"+stepName); err != nil {
			return fmt.Errorf("%s: %w", translator.T("next.error_apply_step", stepName), err)
		}

		if err := manifest.Save(root); err != nil {
			return err
		}

		return stateRepo.Save(s)
	}

	// Case 3: no next step
	if !s.KataFinished {
		s.KataFinished = true
		slog.Info(
			translator.T("next.log_completed"),
			"kata", s.KataName,
		)

		if err := stateRepo.Save(s); err != nil {
			return err
		}

		// Hook: on_kata_finish_hook
		if hookErr := runner.RunOptional("on_kata_finish_hook"); hookErr != nil {
			return fmt.Errorf("on_kata_finish_hook failed: %w", hookErr)
		}

		return nil
	}

	// Case 4: kata already completed
	slog.Info(translator.T("next.log_already_completed"))

	return nil
}
