package main

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
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

			return cmd.Next(root, stateRepo, translator)
		},
	}
}
