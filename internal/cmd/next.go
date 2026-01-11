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
)

func Next(root string, stateRepo *state.Repository, translator i18n.Translator) error {
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
			if errors.Is(err, state.ErrTestsFailing) {
				return errors.New(translator.T("next.error_tests_failing"))
			}
			if errors.Is(err, state.ErrNoNextStep) {
				return errors.New(translator.T("next.error_no_next_step"))
			}
			return err
		}

		stepName := s.Steps[s.CurrentStepIndex]
		slog.Info(translator.T("next.log_advancing"), "step", stepName)

		manifest, err := scaffold.LoadManifest(root)
		if err != nil {
			return err
		}

		copier := &scaffold.Copier{
			Root:       root,
			Manifest:   manifest,
			Translator: translator,
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
		return stateRepo.Save(s)
	}

	// Case 4: kata already completed
	slog.Info(translator.T("next.log_already_completed"))

	return nil
}
