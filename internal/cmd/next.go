package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rjfonseca/kata/internal/fsutil"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/scaffold"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
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
		// --- CREATE CHECKPOINT ---
		// We save the current state (which is Green) before applying the next step.
		checkpointName := "scaffold"
		if s.CurrentStepIndex >= 0 {
			checkpointName = fmt.Sprintf("step-%02d-%s", s.CurrentStepIndex, s.Steps[s.CurrentStepIndex])
		}

		checkpointsDir := filepath.Join(root, ".kata", "checkpoints")
		_ = os.MkdirAll(checkpointsDir, 0o755)
		checkpointPath := filepath.Join(checkpointsDir, checkpointName+".tar.gz")

		// We ignore things that shouldn't be in the checkpoint
		ignore := []string{".git", ".kata", "katas", ".task", "node_modules", "vendor", ".venv"}
		if err := fsutil.CreateArchive(root, checkpointPath, ignore); err != nil {
			slog.Warn("Failed to create checkpoint", "error", err)
		}
		// --- END CHECKPOINT ---

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

		// Hook: on_kata_finish
		if runner, err := taskrunner.New(root); err == nil {
			if err := runner.TryRun("on_kata_finish"); err != nil {
				// We log warning
				slog.Warn(translator.T("next.log_hook_on_kata_finish_failed"), "error", err)
			}
		}

		return stateRepo.Save(s)
	}

	// Case 4: kata already completed
	slog.Info(translator.T("next.log_already_completed"))

	return nil
}
