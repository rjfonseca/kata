package cmd

import (
	"fmt"

	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

func Run(stateRepo *state.Repository, runner taskrunner.Runner) error {
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
		return fmt.Errorf("running task 'test': %w", err)
	}

	return nil
}
