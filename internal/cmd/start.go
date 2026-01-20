package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/rjfonseca/kata/internal/assets"
	"github.com/rjfonseca/kata/internal/config"
	"github.com/rjfonseca/kata/internal/fsutil"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/kata"
	"github.com/rjfonseca/kata/internal/scaffold"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

type StartFlags struct {
	Force bool
}

func Start(workspaceRoot, projectRoot string, stateRepo *state.Repository, kataName string, f StartFlags, translator i18n.Translator) error {
	if stateRepo.Exists() && !f.Force {
		return errors.New(translator.T("start.error_kata_already_started"))
	}

	// ------------------------------------------------------------------
	// 1. Resolve kata source (local catalog wins over embedded)
	// ------------------------------------------------------------------

	localKataDir := filepath.Join(workspaceRoot, "katas", "catalog", kataName)
	// Check if local kata exists, if not check for underscore prefixed version
	if _, err := os.Stat(localKataDir); os.IsNotExist(err) {
		underscoreDir := filepath.Join(workspaceRoot, "katas", "catalog", "_"+kataName)
		if _, err := os.Stat(underscoreDir); err == nil {
			localKataDir = underscoreDir
		}
	}

	var kataFS fs.FS
	var kataBasePath string

	if _, err := os.Stat(localKataDir); err == nil {
		slog.Info(translator.T("start.log_using_local_catalog"), "kata", kataName)
		kataFS = os.DirFS(localKataDir)
		kataBasePath = "."
	} else {
		slog.Info(translator.T("start.log_using_embedded_catalog"), "kata", kataName)
		kataFS = assets.FS
		kataBasePath = filepath.Join("catalog", kataName)

		// Check if it exists in embedded, if not try underscore prefix
		if _, err := fs.Stat(kataFS, kataBasePath); errors.Is(err, fs.ErrNotExist) {
			underscorePath := filepath.Join("catalog", "_"+kataName)
			if _, err := fs.Stat(kataFS, underscorePath); err == nil {
				kataBasePath = underscorePath
			}
		}
	}

	// ------------------------------------------------------------------
	// 2. Copy kata to .kata/catalog/<kata>
	// ------------------------------------------------------------------

	kataDstDir := filepath.Join(projectRoot, ".kata", "catalog", kataName)

	// Clean existing kata if --force was used (Deep Clean)
	if f.Force && stateRepo.Exists() {
		// Load config from Workspace Root to get protected paths
		wsCfg, _ := config.LoadWorkspaceConfig(workspaceRoot)
		protected := wsCfg.GetProtectedPaths()

		// Safety Backup before wipe
		oldState, _ := stateRepo.Load()
		oldName := "unknown"
		if oldState != nil {
			oldName = oldState.KataName
		}

		backupDir := filepath.Join(projectRoot, ".kata", "backups")
		_ = os.MkdirAll(backupDir, 0o755)
		backupPath := filepath.Join(backupDir, fmt.Sprintf("pre-switch-%s-%d.tar.gz", oldName, time.Now().Unix()))

		if err := fsutil.CreateArchive(projectRoot, backupPath, []string{".kata"}); err == nil {
			slog.Info(translator.T("reset.log_safety_backup"), "path", backupPath)
		}

		// Wipe everything not protected in Project Root
		_ = fsutil.CleanDir(projectRoot, protected)
		// Re-create .kata because CleanDir might have removed it if not in defaults (it is in defaults now)
		_ = os.MkdirAll(filepath.Join(projectRoot, ".kata"), 0o755)
	}

	if err := fsutil.CopyDir(kataFS, kataBasePath, kataDstDir); err != nil {
		return fmt.Errorf("%s: %w", translator.T("start.error_copy_failed", kataName), err)
	}

	// ------------------------------------------------------------------
	// 3. Discover steps (from copied catalog)
	// ------------------------------------------------------------------

	stepsDir := filepath.Join(kataDstDir, "steps")

	entries, err := os.ReadDir(stepsDir)
	if err != nil {
		return fmt.Errorf("%s: %w", translator.T("start.error_read_steps_failed", kataName), err)
	}

	var steps []string
	for _, e := range entries {
		if e.IsDir() {
			steps = append(steps, e.Name())
		}
	}

	if len(steps) == 0 {
		return errors.New(translator.T("start.error_no_steps", kataName))
	}

	sort.Strings(steps)

	// ------------------------------------------------------------------
	// Read kata config (runner)
	// ------------------------------------------------------------------

	configPath := filepath.Join(kataDstDir, "config.toml")

	var cfg kata.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return fmt.Errorf("%s: %w", translator.T("start.error_read_config_failed", kataName), err)
	}

	if cfg.Runner.Name == "" {
		return errors.New(translator.T("start.error_no_runner", kataName))
	}

	slog.Info(translator.T("start.log_runner_selected"), "runner", cfg.Runner.Name)

	// ------------------------------------------------------------------
	// Resolve runner source (local wins over embedded)
	// ------------------------------------------------------------------

	runnerName := cfg.Runner.Name
	localRunnerDir := filepath.Join(workspaceRoot, "katas", "runners", runnerName)

	var runnerFS fs.FS
	var runnerBasePath string

	if _, err := os.Stat(localRunnerDir); err == nil {
		slog.Info(translator.T("start.log_using_local_runner"), "runner", runnerName)
		runnerFS = os.DirFS(localRunnerDir)
		runnerBasePath = "."
	} else {
		slog.Info(translator.T("start.log_using_embedded_runner"), "runner", runnerName)
		runnerFS = assets.FS
		runnerBasePath = filepath.Join("runners", runnerName)
	}

	// ------------------------------------------------------------------
	// Copy runner files to project root
	// ------------------------------------------------------------------

	if err := fsutil.CopyDir(
		runnerFS,
		runnerBasePath,
		projectRoot,
	); err != nil {
		return fmt.Errorf("%s: %w", translator.T("start.error_copy_runner_failed", runnerName), err)
	}

	// ------------------------------------------------------------------
	// 4. Initialize state
	// ------------------------------------------------------------------

	s := &state.State{
		KataName:         kataName,
		Runner:           runnerName,
		Steps:            steps,
		CurrentStepIndex: -1,    // Start only with scaffold
		TestPassing:      false, // start in Red
		KataFinished:     false,
		StartedAt:        time.Now(),
	}

	if err := stateRepo.Save(s); err != nil {
		return err
	}

	// ------------------------------------------------------------------
	// 5. Apply scaffold (if present)
	// ------------------------------------------------------------------

	manifest, err := scaffold.LoadManifest(projectRoot)
	if err != nil {
		return err
	}

	copier := &scaffold.Copier{
		Root:         projectRoot,
		Manifest:     manifest,
		TemplateData: scaffold.TemplateData{KataName: kataName},
		Translator:   translator,
	}

	// ------------------------------------------------------------------
	// Apply global scaffold (if present)
	// ------------------------------------------------------------------

	// Local global scaffold check (in Workspace Root)
	localGlobalScaffoldDir := filepath.Join(workspaceRoot, "katas", "scaffold")
	var globalScaffoldFS fs.FS
	var globalScaffoldBasePath string

	if info, err := os.Stat(localGlobalScaffoldDir); err == nil && info.IsDir() {
		// Use local global scaffold
		globalScaffoldFS = os.DirFS(localGlobalScaffoldDir)
		globalScaffoldBasePath = "."
	} else {
		// Use embedded global scaffold
		globalScaffoldFS = assets.FS
		// Use forward slashes for embed.FS, even on Windows
		globalScaffoldBasePath = "scaffold"
	}

	// Apply it
	// We need to check if the path exists in the chosen FS.
	// For embedded, we check if scaffold exists.
	// For local, we already checked directory existence.
	shouldApplyGlobal := true
	if _, err := fs.Stat(globalScaffoldFS, globalScaffoldBasePath); err != nil {
		shouldApplyGlobal = false
	}

	if shouldApplyGlobal {
		// We need a sub-fs for the apply
		sub, err := fs.Sub(globalScaffoldFS, globalScaffoldBasePath)
		if err == nil {
			if err := copier.Apply(sub, "scaffold:global"); err != nil {
				// We don't fail here to keep backward compatibility or if scaffold is missing
				// But we log it
				slog.Warn("Failed to apply global scaffold", "error", err)
			}
		}
	}

	scaffoldDir := filepath.Join(kataDstDir, "scaffold")
	if info, err := os.Stat(scaffoldDir); err == nil && info.IsDir() {
		slog.Info(translator.T("start.log_scaffold_applied"))
		if err := copier.Apply(os.DirFS(scaffoldDir), "scaffold"); err != nil {
			return fmt.Errorf("%s: %w", translator.T("start.error_scaffold_failed", kataName), err)
		}
	}

	if err := manifest.Save(projectRoot); err != nil {
		return err
	}

	// ------------------------------------------------------------------
	// 6. Execute Hook: on_kata_start
	// ------------------------------------------------------------------
	// We create a runner here to execute the hook.
	// If the Taskfile is invalid, this might fail, which is acceptable.
	if runner, err := taskrunner.New(projectRoot); err == nil {
		if err := runner.TryRun("on_kata_start"); err != nil {
			return fmt.Errorf("%s: %w", translator.T("start.error_hook_on_kata_start_failed"), err)
		}
	} else {
		// Just log, don't fail the start if we can't load the runner (maybe no Taskfile yet?)
		// But we just copied it.
		slog.Debug("failed to create task runner for hooks", "error", err)
	}

	slog.Info(translator.T("start.log_kata_started"), "kata", kataName)
	return nil

}
