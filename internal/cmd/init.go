package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rjfonseca/kata/internal/assets"
	"github.com/rjfonseca/kata/internal/fsutil"
	"github.com/rjfonseca/kata/internal/i18n"
)

func Init(root string, translator i18n.Translator) error {
	katasDir := filepath.Join(root, "katas")
	catalogDir := filepath.Join(katasDir, "catalog")
	scaffoldDir := filepath.Join(katasDir, "scaffold")
	runnersDir := filepath.Join(katasDir, "runners")
	i18nDir := filepath.Join(katasDir, "i18n")

	if _, err := os.Stat(katasDir); err == nil {
		slog.Info(translator.T("init.log_dir_exists"))
		return nil
	}

	slog.Info(translator.T("init.log_creating_dir"), "path", katasDir)
	if err := os.MkdirAll(katasDir, 0o755); err != nil {
		return err
	}

	slog.Info(translator.T("init.log_copying_catalog"))
	if err := fsutil.CopyDir(assets.FS, "catalog", catalogDir); err != nil {
		return err
	}

	slog.Info(translator.T("init.log_copying_scaffold"))
	if err := fsutil.CopyDir(assets.FS, "scaffold", scaffoldDir); err != nil {
		return err
	}

	slog.Info(translator.T("init.log_copying_runners"))
	if err := fsutil.CopyDir(assets.FS, "runners", runnersDir); err != nil {
		return err
	}

	slog.Info(translator.T("init.log_copying_i18n"))
	if err := fsutil.CopyDir(assets.FS, "i18n", i18nDir); err != nil {
		return err
	}

	// Copy kata.toml to root
	slog.Info(translator.T("init.log_creating_config"))
	configContent, err := assets.FS.ReadFile("kata.toml")
	if err == nil {
		_ = os.WriteFile(filepath.Join(root, "kata.toml"), configContent, 0o644)
	}

	slog.Info(translator.T("init.log_repo_initialized"))
	return nil
}
