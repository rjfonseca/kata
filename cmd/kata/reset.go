package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
)

func resetCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:  "reset",
		Usage: translator.T("reset.usage"),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "hard",
				Usage: translator.T("reset.flag_hard_usage"),
			},
			&cli.IntFlag{
				Name:  "to-step",
				Value: -2, // Using -2 as unset since -1 is scaffold
				Usage: translator.T("reset.flag_to_step_usage"),
			},
		},
		Action: func(c *cli.Context) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}

			stateRepo := state.NewRepository(root)

			f := cmd.ResetFlags{
				Hard:   c.Bool("hard"),
				ToStep: c.Int("to-step"),
			}

			// If to-step was not provided, but hard was, it's already handled.
			// If neither was provided, cmd.Reset uses CurrentStepIndex.
			if f.ToStep == -2 {
				f.ToStep = -10 // Sentinel for "not provided"
			}

			if err := cmd.Reset(root, stateRepo, f, translator); err != nil {
				return err
			}

			return nil
		},
	}
}
