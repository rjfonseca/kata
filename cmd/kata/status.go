package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
)

func statusCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: translator.T("status.usage"),
		Action: func(c *cli.Context) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}

			stateRepo := state.NewRepository(root)
			st, err := stateRepo.Load()
			if err != nil {
				return fmt.Errorf("%s: %w", translator.T("status.error_read_state"), err)
			}

			// Styles
			labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Width(16)
			valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true)
			successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
			failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
			infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)

			printRow := func(label, value string, style lipgloss.Style) {
				fmt.Printf("%s %s\n", labelStyle.Render(label), style.Render(value))
			}

			// Kata Name
			printRow(translator.T("status.field_kata"), st.KataName, infoStyle)

			// Current Step
			currentStep := "?"
			if st.CurrentStepIndex >= 0 && st.CurrentStepIndex < len(st.Steps) {
				currentStep = st.Steps[st.CurrentStepIndex]
			}
			printRow(translator.T("status.field_step"), currentStep, valueStyle)

			// Passing Status
			passing := translator.T("status.value_no")
			style := failStyle
			if st.TestPassing {
				passing = translator.T("status.value_yes")
				style = successStyle
			}
			printRow(translator.T("status.field_passing"), passing, style)

			// Completed Status
			completed := translator.T("status.value_no")
			style = valueStyle
			if st.KataFinished {
				completed = translator.T("status.value_yes")
				style = successStyle
			}
			printRow(translator.T("status.field_completed"), completed, style)

			return nil
		},
	}
}
