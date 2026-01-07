package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/interactive"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

// runCommand creates the `kata run` command.
func runCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: translator.T("run.usage"),
		Action: func(ctx *cli.Context) error {
			root, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("%s: %w", translator.T("run.error_read_cwd"), err)
			}

			stateRepo := state.NewRepository(root)

			runner, err := taskrunner.New(root)
			if err != nil {
				return fmt.Errorf("%s: %w", translator.T("run.error_create_executor"), err)
			}
			err = cmd.Run(stateRepo, runner, translator)
			if isNonInteractive(ctx) {
				return err
			}

			return interactive.Run(ctx.Context, interactive.Options{
				LoadState: stateRepo.Load,
				Run: func() error {
					return cmd.Run(stateRepo, runner, translator)
				},
				Next: func() error {
					return cmd.Next(root, stateRepo, translator)
				},
			}, translator)
		},
	}
}
