package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/interactive"
	"github.com/rjfonseca/kata/internal/kata"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
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
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			stateRepo := state.NewRepository(root)

			kataName := c.Args().First()
			f := cmd.StartFlags{
				Force: c.Bool("force"),
			}

			// Non-interactive path
			if kataName == "" {
				if isNonInteractive(c) {
					return errors.New(translator.T("start.error_kata_name_required"))
				}

				kataRepo := kata.NewRepository(root)
				katas, err := kataRepo.ListKatas()
				if err != nil {
					return fmt.Errorf("%s: %w", translator.T("start.error_list_katas"), err)
				}

				if len(katas) == 0 {
					fmt.Println(translator.T("start.log_no_katas_found"))
					return nil
				}

				kataName, err = interactive.Select(translator.T("start.prompt_select_kata"), katas)
				if err != nil {
					return fmt.Errorf("%s: %w", translator.T("start.error_select_kata"), err)
				}
				if kataName == "" {
					return nil // User cancelled
				}
			}

			if err := cmd.Start(root, stateRepo, kataName, f, translator); err != nil {
				return err
			}

			if !isNonInteractive(c) {
				runner, err := taskrunner.New(root)
				if err != nil {
					return fmt.Errorf("%s: %w", translator.T("run.error_create_executor"), err)
				}

				// First run after start
				_ = cmd.Run(stateRepo, runner, translator)

				return interactive.Run(c.Context, interactive.Options{
					LoadState: stateRepo.Load,
					Run: func() error {
						return cmd.Run(stateRepo, runner, translator)
					},
					Next: func() error {
						return cmd.Next(root, stateRepo, translator)
					},
				})
			}

			return nil
		},
	}
}
