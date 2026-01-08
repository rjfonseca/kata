package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

func nextCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:  "next",
		Usage: translator.T("next.usage"),
		Action: func(c *cli.Context) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}

			stateRepo := state.NewRepository(root)
			if !stateRepo.Exists() {
				return errors.New(translator.T("next.error_not_started"))
			}

			runner, err := taskrunner.New(root)
			if err != nil {
				return fmt.Errorf("%s: %w", translator.T("run.error_create_executor"), err)
			}

			return cmd.Next(root, stateRepo, runner, translator)
		},
	}
}
