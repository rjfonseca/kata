package main

import (
	"fmt"
	"os"
	"os/exec"

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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   translator.T("run.flag_watch_usage"),
			},
		},
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

			// Always run once
			err = cmd.Run(stateRepo, runner, translator)

			watch := ctx.Bool("watch")

			// If not interactive and not watching, return the result of the run immediately
			if !isInteractive(ctx) && !watch {
				return err
			}

			exe, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %w", err)
			}

			// If interactive or watching, enter the loop
			return interactive.Run(ctx.Context, interactive.Options{
				LoadState: stateRepo.Load,
				RunCmd: func() *exec.Cmd {
					return exec.Command(exe, "run")
				},
				NextCmd: func() *exec.Cmd {
					return exec.Command(exe, "next")
				},
				ListTasks: runner.ListTasks,
				RunTaskCmd: func(name string) *exec.Cmd {
					return exec.Command(exe, "task", name)
				},
				Watch: watch,
			})
		},
	}
}
