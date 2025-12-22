package main

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/state"
)

func nextCommand() *cli.Command {
	return &cli.Command{
		Name:  "next",
		Usage: "Advance to the next kata step",
		Action: func(c *cli.Context) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}

			stateRepo := state.NewRepository(root)
			if !stateRepo.Exists() {
				return errors.New("no kata started (run 'kata start <name>' first)")
			}

			return cmd.Next(root, stateRepo)
		},
	}
}
