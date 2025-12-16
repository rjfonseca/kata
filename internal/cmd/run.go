package cmd

import (
	"fmt"
	"log/slog"

	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

func Run(stateRepo *state.Repository, runner taskrunner.Runner, translator i18n.Translator) error {
	st, err := stateRepo.Load()
	if err != nil {
		return err
	}

	// Hook: before_run
	// If this fails, we abort the run.
	if err := runner.TryRun("before_run"); err != nil {
		return fmt.Errorf("%s: %w", translator.T("run.error_hook_before_run_failed"), err)
	}

	err = runner.Run("test")
	st.MarkRunResult(err == nil)

	// Hooks: on_run_success / on_run_fail
	if err == nil {
		if hookErr := runner.TryRun("on_run_success"); hookErr != nil {
			// We log but don't fail the command itself, or maybe we should?
			// User preference usually: if hook fails, it's annoying.
			// But the test passed.
			slog.Warn(translator.T("run.log_hook_on_run_success_failed"), "error", hookErr)
		}
	} else {
		if hookErr := runner.TryRun("on_run_fail"); hookErr != nil {
			slog.Warn(translator.T("run.log_hook_on_run_fail_failed"), "error", hookErr)
		}
	}

	// Hook: on_run (always)
	if hookErr := runner.TryRun("on_run"); hookErr != nil {
		slog.Warn(translator.T("run.log_hook_on_run_failed"), "error", hookErr)
	}

	if saveErr := stateRepo.Save(st); saveErr != nil {
		return saveErr
	}

	if err != nil {
		return fmt.Errorf("%s: %w", translator.T("run.error_run_task"), err)
	}

	return nil
}
