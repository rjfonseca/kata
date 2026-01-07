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

	err = runner.Run("test")
	st.MarkRunResult(err == nil)

	if saveErr := stateRepo.Save(st); saveErr != nil {
		return saveErr
	}

	if err != nil {
		return fmt.Errorf("%s: %w", translator.T("run.error_run_task"), err)
	}

	return nil
}
