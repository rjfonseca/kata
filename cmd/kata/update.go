package main

import (
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
				slog.Info(translator.T("update.error_dev_build"))
				return nil
			}

			slog.Info(translator.T("update.log_checking"))

			latest, found, err := selfupdate.DetectLatest("rjfonseca/kata")
			if err != nil {
				slog.Error(translator.T("update.error_check_failed"), "error", err)
				return cli.Exit(translator.T("update.error_check_failed"), 1)
			}

			if !found {
				slog.Info(translator.T("update.log_already_latest", version))
				return nil
			}

			v, err := semver.Parse(version)
			if err != nil {
				// If current version is not semver, we assume it's older or broken, but usually we just log error
				slog.Warn("failed to parse current version", "version", version, "error", err)
				// Proceeding might be risky if version is completely wrong, but let's try strict check
			}

			if latest.Version.LE(v) {
				slog.Info(translator.T("update.log_already_latest", version))
				return nil
			}

			slog.Info(translator.T("update.log_update_found", latest.Version.String()))
			slog.Info(translator.T("update.log_updating", latest.Version.String()))

			exe, err := os.Executable()
			if err != nil {
				slog.Error("failed to locate executable", "error", err)
				return cli.Exit(translator.T("update.error_update_failed"), 1)
			}

			if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
				slog.Error(translator.T("update.error_update_failed"), "error", err)
				return cli.Exit(translator.T("update.error_update_failed"), 1)
			}

			slog.Info(translator.T("update.log_success", latest.Version.String()))
			return nil
		},
	}
}
