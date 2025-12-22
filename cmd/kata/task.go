package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/i18n"
)

func taskCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:            "task",
		Aliases:         []string{"t"},
		Usage:           translator.T("task.usage"),
		SkipFlagParsing: true, // Allow flags to be passed to the task (future proofing)
		Action: func(c *cli.Context) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			return cmd.Task(root, c.Args().Slice(), translator)
		},
	}
}
