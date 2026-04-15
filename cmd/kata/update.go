package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/urfave/cli/v2"
)

func updateCommand(translator i18n.Translator) *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: translator.T("update.usage"),
		Action: func(c *cli.Context) error {
			if version == "dev" {
				slog.Info(translator.T("update.log_dev_build"))
				return nil
			}

			slog.Info(translator.T("update.log_checking"))

			latest, found, err := selfupdate.DetectLatest("rjfonseca/kata")
			if err != nil {
				return fmt.Errorf("%s: %w", translator.T("update.error_check_failed"), err)
			}

			if !found {
				slog.Info(translator.T("update.log_already_latest", version))
				return nil
			}

			v, err := semver.ParseTolerant(version)
			if err != nil {
				return fmt.Errorf("%s: %w", translator.T("update.error_parse_version", version), err)
			}

			if latest.Version.LE(v) {
				slog.Info(translator.T("update.log_already_latest", version))
				return nil
			}

			slog.Info(translator.T("update.log_update_found", latest.Version.String()))
			slog.Info(translator.T("update.log_updating", latest.Version.String()))

			exe, err := os.Executable()
			if err != nil {
				return fmt.Errorf("%s: %w", translator.T("update.error_locate_executable"), err)
			}

			if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
				return fmt.Errorf("%s: %w", translator.T("update.error_update_failed"), err)
			}

			slog.Info(translator.T("update.log_success", latest.Version.String()))
			return nil
		},
	}
}
