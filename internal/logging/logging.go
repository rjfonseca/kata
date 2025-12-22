package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// Init initializes the global slog logger for CLI usage.
func Init(level slog.Level) {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		TimeFormat: "", // disable timestamps
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Drop time attribute for CLI output
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
