package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rjfonseca/kata/internal/scaffold"
	"github.com/rjfonseca/kata/internal/state"
)

func Next(root string, stateRepo *state.Repository) error {
	s, err := stateRepo.Load()
	if err != nil {
		return err
	}

	// Case 1: tests are failing (Red)
	if !s.TestPassing {
		return errors.New("tests are failing, make them pass before advancing")
	}

	// Case 2: there is a next step to advance to
	if s.HasNextStep() {
		if err := s.Next(); err != nil {
			return err
		}

		stepName := s.Steps[s.CurrentStepIndex]
		slog.Info("advancing to next step", "step", stepName)

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
			return fmt.Errorf("failed to apply step %q: %w", stepName, err)
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
			"kata completed ð",
			"kata", s.KataName,
		)
		return stateRepo.Save(s)
	}

	// Case 4: kata already completed
	slog.Info(
		"kata already completed! You can keep refactoring and learning",
	)

	return nil
}
