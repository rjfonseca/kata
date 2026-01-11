package main

import (
	"log/slog"

	"github.com/urfave/cli/v2"

	"github.com/rjfonseca/kata/internal/config"
	"github.com/rjfonseca/kata/internal/i18n"
)

func configCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: translator.T("config.usage"),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "set-lang",
				Usage: translator.T("config.flag_set_lang_usage"),
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
				slog.Info(translator.T("config.log_lang_updated"), "lang", lang)
				return nil
			}

			slog.Debug("Configuration", "ui.lang", cfg.UI.Lang)
			return nil
		},
	}
}
