package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/rjfonseca/kata/internal/config"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/i18n/loader"
	"github.com/rjfonseca/kata/internal/logging"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "kata",
		Usage: "Create and run code katas",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "lang",
				Usage: "Language for messages (e.g. en-US, pt-BR)",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "enable verbose (debug) logging",
			},
			&cli.BoolFlag{
				Name:  "quiet",
				Usage: "suppress informational logs",
			},
			&cli.BoolFlag{
				Name:  "non-interactive",
				Usage: "disables interactive mode. Use this for scripting, ci or AI agents",
			},
		},
		Before: func(ctx *cli.Context) error {
			level, err := resolveLogLevel(ctx)
			if err != nil {
				return err
			}
			logging.Init(level)

			// --- I18n Initialization ---
			langFlag := ctx.String("lang")
			langTag := config.ResolveLanguage(langFlag)

			// Assume current directory is the project root for file loading.
			projectRoot := "."
			embeddedLoader := loader.NewEmbeddedLoader()
			fileLoader := loader.NewFileLoader(projectRoot)
			compositeLoader := loader.NewCompositeLoader(embeddedLoader, fileLoader)

			translator, err := i18n.NewDefaultTranslator(compositeLoader, langTag)
			if err != nil {
				return fmt.Errorf("failed to initialize translator: %w", err)
			}

			if ctx.App.Metadata == nil {
				ctx.App.Metadata = make(map[string]interface{})
			}
			ctx.App.Metadata["translator"] = translator
			// --- End I18n Initialization ---

			return nil
		},
		Commands: []*cli.Command{
			initCommand(),
			configCommand(),
			startCommand(),
			nextCommand(),
			runCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("application error", "err", err)
		os.Exit(1)
	}
}

func resolveLogLevel(ctx *cli.Context) (slog.Level, error) {
	verbose := ctx.Bool("verbose")
	quiet := ctx.Bool("quiet")

	if verbose && quiet {
		return slog.LevelInfo, fmt.Errorf("--verbose and --quiet cannot be used together")
	}

	switch {
	case verbose:
		return slog.LevelDebug, nil
	case quiet:
		return slog.LevelWarn, nil
	default:
		return slog.LevelInfo, nil
	}
}

func isNonInteractive(ctx *cli.Context) bool {
	return ctx.Bool("non-interactive")
}
