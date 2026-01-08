package cmd

import (
	"fmt"

	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

func Run(stateRepo *state.Repository, runner taskrunner.Runner, translator i18n.Translator) error {
	st, err := stateRepo.Load()
	if err != nil {
		return err
	}

	// Hook: pre_run_hook
	if hookErr := runner.RunOptional("pre_run_hook"); hookErr != nil {
		return fmt.Errorf("pre_run_hook failed: %w", hookErr)
	}

	err = runner.Run("test")
	st.MarkRunResult(err == nil)

	if saveErr := stateRepo.Save(st); saveErr != nil {
		return saveErr
	}

	if err != nil {
		// Hook: on_failure_hook
		if hookErr := runner.RunOptional("on_failure_hook"); hookErr != nil {
			return fmt.Errorf("on_failure_hook failed: %w", hookErr)
		}
		return fmt.Errorf("%s: %w", translator.T("run.error_run_task"), err)
	}

	// Hook: on_success_hook
	if hookErr := runner.RunOptional("on_success_hook"); hookErr != nil {
		return fmt.Errorf("on_success_hook failed: %w", hookErr)
	}

	return nil
}
