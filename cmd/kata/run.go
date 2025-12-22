package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/interactive"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

// runCommand creates the `kata run` command.
func runCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Run the kata tests",
		Action: func(ctx *cli.Context) error {

			root, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("reading working directory: %w", err)
			}

			stateRepo := state.NewRepository(root)

			runner, err := taskrunner.New()
			if err != nil {
				return fmt.Errorf("creating task executor: %w", err)
			}
			err = cmd.Run(stateRepo, runner)
			if isNonInteractive(ctx) {
				return err
			}

			return interactive.Run(ctx.Context, interactive.Options{
				LoadState: stateRepo.Load,
				Run: func() error {
					return cmd.Run(stateRepo, runner)
				},
				Next: func() error {
					return cmd.Next(root, stateRepo)
				},
			})
		},
	}
}
