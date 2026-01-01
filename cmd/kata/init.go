package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
	"github.com/rjfonseca/kata/internal/i18n"
)

// initCommand initializes a kata catalog repository in the current directory.
func initCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: translator.T("init.usage"),
		Action: func(c *cli.Context) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			return cmd.Init(root, translator)
		},
	}
}
