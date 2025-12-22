package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/config"
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
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   translator.T("start.flag_dir_usage"),
			},
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   translator.T("run.flag_watch_usage"),
			},
		},
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			// 1. Resolve Workspace Root (for catalog lookup)
			workspaceRoot, err := config.FindWorkspaceRoot(cwd)
			if err != nil {
				if errors.Is(err, config.ErrNoWorkspace) {
					// Fallback: use CWD (embedded catalog will be used if local katas not found)
					workspaceRoot = cwd
				} else {
					return err
				}
			}

			// 2. Resolve Project Root (where kata will be created)
			projectRoot := cwd
			if dir := c.String("dir"); dir != "" {
				projectRoot = filepath.Join(cwd, dir)
				if err := os.MkdirAll(projectRoot, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", projectRoot, err)
				}
				slog.Info("Created directory", "path", projectRoot)
			}

			stateRepo := state.NewRepository(projectRoot)

			kataName := c.Args().First()
			f := cmd.StartFlags{
				Force: c.Bool("force"),
			}

			// If no name provided, allow selection (regardless of -i flag, as this is argument input)
			if kataName == "" {
				kataRepo := kata.NewRepository(workspaceRoot)
				katas, err := kataRepo.ListKatas()
				if err != nil {
					return fmt.Errorf("%s: %w", translator.T("start.error_list_katas"), err)
				}

				if len(katas) == 0 {
					slog.Info(translator.T("start.log_no_katas_found"))
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

			if err := cmd.Start(workspaceRoot, projectRoot, stateRepo, kataName, f, translator); err != nil {
				return err
			}

			watch := c.Bool("watch")

			// Only enter the interactive loop if explicitly requested or watch is enabled
			if isInteractive(c) || watch {
				runner, err := taskrunner.New(projectRoot)
				if err != nil {
					return fmt.Errorf("%s: %w", translator.T("run.error_create_executor"), err)
				}

				// First run after start
				_ = cmd.Run(stateRepo, runner, translator)

				exe, err := os.Executable()
				if err != nil {
					return fmt.Errorf("failed to get executable path: %w", err)
				}

				return interactive.Run(c.Context, interactive.Options{
					LoadState: stateRepo.Load,
					RunCmd: func() *exec.Cmd {
						c := exec.Command(exe, "run")
						c.Dir = projectRoot
						return c
					},
					NextCmd: func() *exec.Cmd {
						c := exec.Command(exe, "next")
						c.Dir = projectRoot
						return c
					},
					ListTasks: runner.ListTasks,
					RunTaskCmd: func(name string) *exec.Cmd {
						c := exec.Command(exe, "task", name)
						c.Dir = projectRoot
						return c
					},
					Watch: watch,
				})
			}

			return nil
		},
	}
}
