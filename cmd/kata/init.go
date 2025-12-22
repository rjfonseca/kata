package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/cmd"
)

// initCommand initializes a kata catalog repository in the current directory.
func initCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize a kata catalog repository",
		Action: func(c *cli.Context) error {
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			return cmd.Init(root)
		},
	}
}
