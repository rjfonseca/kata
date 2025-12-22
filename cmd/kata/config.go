package main

import (
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/config"
)

func configCommand() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Manage kata configuration",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "set-lang",
				Usage: "Set global language (e.g. pt-BR)",
			},
		},
		Action: func(c *cli.Context) error {
			repo, err := config.NewRepository()
			if err != nil {
				return err
			}

			cfg, err := repo.Load()
			if err != nil {
				return err
			}

			if lang := c.String("set-lang"); lang != "" {
				cfg.UI.Lang = lang
				if err := repo.Save(cfg); err != nil {
					return err
				}
				slog.Info("language updated", "lang", lang)
				return nil
			}

			fmt.Printf("ui.lang = %s\n", cfg.UI.Lang)
			return nil
		},
	}
}
