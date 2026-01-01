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
func startCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:      "start",
		Usage:     translator.T("start.usage"),
		ArgsUsage: translator.T("start.args_usage"),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: translator.T("start.flag_force_usage"),
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return errors.New(translator.T("start.error_kata_name_required"))
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
