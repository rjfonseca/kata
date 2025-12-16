package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/rjfonseca/kata/internal/config"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/i18n/loader"
	"github.com/rjfonseca/kata/internal/logging"
	"github.com/urfave/cli/v2"
)

func main() {
	// --- Pre-flight and i18n Initialization ---
	logging.Init(slog.LevelWarn)

	root, err := os.Getwd()
	if err != nil {
		slog.Error("Unable to detect current working directory", "error", err)
		os.Exit(1)
	}
	langCode := preflightLang(os.Args)
	langTag := config.ResolveLanguage(langCode)
	translator, err := i18n.NewDefaultTranslator(
		loader.NewCompositeLoader(
			loader.NewEmbeddedLoader(),
			loader.NewFileLoader(root),
		),
		langTag,
	)
	if err != nil {
		// Cannot use slog yet, and cannot use i18n.
		fmt.Fprintf(os.Stderr, "failed to initialize translator: %v\n", err)
		os.Exit(1)
	}

	// --- CLI App Definition ---
	app := &cli.App{
		Name:                 "kata",
		Usage:                translator.T("cli.usage"),
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "lang",
				Usage: translator.T("cli.flag_lang_usage"),
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: translator.T("cli.flag_verbose_usage"),
			},
			&cli.BoolFlag{
				Name:  "quiet",
				Usage: translator.T("cli.flag_quiet_usage"),
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   translator.T("cli.flag_interactive_usage"),
			},
		},
		Before: func(ctx *cli.Context) error {
			// Set up logging
			level, err := resolveLogLevel(ctx)
			if err != nil {
				return err
			}
			logging.Init(level)

			// Store the translator in the context for sub-commands
			if ctx.App.Metadata == nil {
				ctx.App.Metadata = make(map[string]interface{})
			}
			ctx.App.Metadata["translator"] = translator
			return nil
		},
		Commands: []*cli.Command{
			initCommand(translator),
			configCommand(translator),
			startCommand(translator),
			nextCommand(translator),
			runCommand(translator),
			resetCommand(translator),
			listCommand(translator),
			statusCommand(translator),
			completionCommand(translator),
			taskCommand(translator),
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("application error", "err", err)
		os.Exit(1)
	}
}

// preflightLang manually parses os.Args to find the --lang flag value
// before the main CLI app parsing begins. This allows i18n to be
// initialized early enough to translate help text.
func preflightLang(args []string) string {
	for i, arg := range args {
		if (arg == "--lang" || arg == "-l") && i+1 < len(args) {
			return args[i+1]
		}
		if lang, found := strings.CutPrefix(arg, "--lang="); found {
			return lang
		}
	}
	return ""
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

func isInteractive(ctx *cli.Context) bool {
	return ctx.Bool("interactive")
}
