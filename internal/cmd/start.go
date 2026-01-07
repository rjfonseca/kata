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
	"github.com/rjfonseca/kata/internal/fsutil"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/kata"
	"github.com/rjfonseca/kata/internal/scaffold"
	"github.com/rjfonseca/kata/internal/state"
)

type StartFlags struct {
	Force bool
}

func Start(root string, stateRepo *state.Repository, kataName string, f StartFlags, translator i18n.Translator) error {
	if stateRepo.Exists() && !f.Force {
		return errors.New(translator.T("start.error_kata_already_started"))
	}

	// ------------------------------------------------------------------
	// 1. Resolve kata source (local catalog wins over embedded)
	// ------------------------------------------------------------------

	localKataDir := filepath.Join(root, "katas", "catalog", kataName)

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
	}

	// ------------------------------------------------------------------
	// 2. Copy kata to .kata/catalog/<kata>
	// ------------------------------------------------------------------

	kataDstDir := filepath.Join(root, ".kata", "catalog", kataName)

	// Clean existing .kata if --force was used
	if f.Force {
		_ = os.RemoveAll(filepath.Join(root, ".kata"))
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
	localRunnerDir := filepath.Join(root, "katas", "runners", runnerName)

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
		root,
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
		CurrentStepIndex: 0,
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

	manifest, err := scaffold.LoadManifest(root)
	if err != nil {
		return err
	}

	copier := &scaffold.Copier{
		Root:         root,
		Manifest:     manifest,
		TemplateData: scaffold.TemplateData{KataName: kataName},
	}

	// ------------------------------------------------------------------
	// Apply global scaffold (if present)
	// ------------------------------------------------------------------

	// Local global scaffold check
	localGlobalScaffoldDir := filepath.Join(root, "katas", "catalog", "scaffold")
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
		globalScaffoldBasePath = "catalog/scaffold"
	}

	// Apply it
	// We need to check if the path exists in the chosen FS.
	// For embedded, we check if catalog/scaffold exists.
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

	// ------------------------------------------------------------------
	// 6. Apply first step
	// ------------------------------------------------------------------

	firstStep := steps[0]
	firstStepDir := filepath.Join(kataDstDir, "steps", firstStep)

	slog.Info(translator.T("start.log_first_step_applied"), "step", firstStep)

	if err := copier.Apply(
		os.DirFS(firstStepDir),
		"step:"+firstStep,
	); err != nil {
		return fmt.Errorf("%s: %w", translator.T("start.error_apply_first_step_failed", firstStep, kataName), err)
	}

	if err := manifest.Save(root); err != nil {
		return err
	}

	slog.Info(translator.T("start.log_kata_started"), "kata", kataName)
	return nil

}
