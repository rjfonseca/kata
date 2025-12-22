package cmd

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rjfonseca/kata/internal/assets"
	"github.com/rjfonseca/kata/internal/fsutil"
)

func Init(root string) error {
	katasDir := filepath.Join(root, "katas")
	catalogDir := filepath.Join(katasDir, "catalog")
	runnersDir := filepath.Join(katasDir, "runners")

	if _, err := os.Stat(katasDir); err == nil {
		return errors.New("katas directory already exists")
	}

	slog.Info("creating katas directory", "path", katasDir)
	if err := os.MkdirAll(katasDir, 0o755); err != nil {
		return err
	}

	slog.Info("copying embedded catalog")
	if err := fsutil.CopyDir(assets.FS, "catalog", catalogDir); err != nil {
		return err
	}

	slog.Info("copying embedded runners")
	if err := fsutil.CopyDir(assets.FS, "runners", runnersDir); err != nil {
		return err
	}

	slog.Info("kata repository initialized")
	return nil
}
