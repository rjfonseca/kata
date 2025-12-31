package main

import (
	"errors"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
)

// startCommand starts a kata execution.
// It resolves the kata from the local catalog or the embedded catalog,
// copies it to .kata, applies the scaffold and the first step,
// and initializes the kata state.
func startCommand() *cli.Command {
	return &cli.Command{
		Name:      "start",
		Usage:     "Start a kata",
		ArgsUsage: "<kata-name>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Overwrite existing kata state",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return errors.New("kata name is required")
			}

			translator, ok := c.App.Metadata["translator"].(i18n.Translator)
			if !ok {
				// This should not happen if the Before hook is set up correctly
				return errors.New("translator not found in context")
			}

			kataName := c.Args().First()

			root, err := os.Getwd()
			if err != nil {
				return err
			}

			stateRepo := state.NewRepository(root)

			f := cmd.StartFlags{
				Force: c.Bool("force"),
			}

			return cmd.Start(root, stateRepo, kataName, f, translator)
		},
	}
}
