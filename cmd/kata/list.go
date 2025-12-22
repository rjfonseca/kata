package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/config"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/kata"
)

func listCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   translator.T("list.usage"),
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			// Resolve Workspace Root (for catalog lookup)
			workspaceRoot, err := config.FindWorkspaceRoot(cwd)
			if err != nil {
				if errors.Is(err, config.ErrNoWorkspace) {
					// Fallback: use CWD
					workspaceRoot = cwd
				} else {
					return err
				}
			}

			repo := kata.NewRepository(workspaceRoot)
			katas, err := repo.ListKatas()
			if err != nil {
				return fmt.Errorf("%s: %w", translator.T("list.error_fetch"), err)
			}

			if len(katas) == 0 {
				fmt.Println(translator.T("start.log_no_katas_found"))
				return nil
			}

			titleStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Bold(true).
				MarginBottom(1)

			itemStyle := lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("252"))

			fmt.Println(titleStyle.Render("Available Katas:"))

			for _, k := range katas {
				fmt.Println(itemStyle.Render("• " + k))
			}

			return nil
		},
	}
}
