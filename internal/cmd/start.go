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
		return errors.New("kata already started (use --force to overwrite)")
	}

	// ------------------------------------------------------------------
	// 1. Resolve kata source (local catalog wins over embedded)
	// ------------------------------------------------------------------

	localKataDir := filepath.Join(root, "katas", "catalog", kataName)

	var kataFS fs.FS
	var kataBasePath string

	if _, err := os.Stat(localKataDir); err == nil {
		slog.Info("using local kata catalog", "kata", kataName)
		kataFS = os.DirFS(localKataDir)
		kataBasePath = "."
	} else {
		slog.Info("using embedded kata catalog", "kata", kataName)
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
		return fmt.Errorf(
			"failed to copy kata '%s' to .kata: %w",
			kataName,
			err,
		)
	}

	// ------------------------------------------------------------------
	// 3. Discover steps (from copied catalog)
	// ------------------------------------------------------------------

	stepsDir := filepath.Join(kataDstDir, "steps")

	entries, err := os.ReadDir(stepsDir)
	if err != nil {
		return fmt.Errorf(
			"failed to read steps for kata '%s': %w",
			kataName,
			err,
		)
	}

	var steps []string
	for _, e := range entries {
		if e.IsDir() {
			steps = append(steps, e.Name())
		}
	}

	if len(steps) == 0 {
		return fmt.Errorf(
			"kata '%s' has no steps",
			kataName,
		)
	}

	sort.Strings(steps)

	// ------------------------------------------------------------------
	// Read kata config (runner)
	// ------------------------------------------------------------------

	configPath := filepath.Join(kataDstDir, "config.toml")

	var cfg kata.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return fmt.Errorf(
			"failed to read config.toml for kata '%s': %w",
			kataName,
			err,
		)
	}

	if cfg.Runner == "" {
		return fmt.Errorf(
			"kata '%s' does not define a runner in config.toml",
			kataName,
		)
	}

	slog.Info("kata runner selected", "runner", cfg.Runner)

	// ------------------------------------------------------------------
	// Resolve runner source (local wins over embedded)
	// ------------------------------------------------------------------

	runnerName := cfg.Runner
	localRunnerDir := filepath.Join(root, "katas", "runners", runnerName)

	var runnerFS fs.FS
	var runnerBasePath string

	if _, err := os.Stat(localRunnerDir); err == nil {
		slog.Info("using local runner", "runner", runnerName)
		runnerFS = os.DirFS(localRunnerDir)
		runnerBasePath = "."
	} else {
		slog.Info("using embedded runner", "runner", runnerName)
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
		return fmt.Errorf(
			"failed to copy runner '%s' to project root: %w",
			runnerName,
			err,
		)
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
		Root:     root,
		Manifest: manifest,
	}

	scaffoldDir := filepath.Join(kataDstDir, "scaffold")
	if info, err := os.Stat(scaffoldDir); err == nil && info.IsDir() {
		slog.Info("applying kata scaffold")
		if err := copier.Apply(os.DirFS(scaffoldDir), "scaffold"); err != nil {
			return fmt.Errorf(
				"failed to apply scaffold for kata '%s': %w",
				kataName,
				err,
			)
		}
	}

	// ------------------------------------------------------------------
	// 6. Apply first step
	// ------------------------------------------------------------------

	firstStep := steps[0]
	firstStepDir := filepath.Join(kataDstDir, "steps", firstStep)

	slog.Info("applying first kata step", "step", firstStep)

	if err := copier.Apply(
		os.DirFS(firstStepDir),
		"step:"+firstStep,
	); err != nil {
		return fmt.Errorf(
			"failed to apply first step '%s' for kata '%s': %w",
			firstStep,
			kataName,
			err,
		)
	}

	if err := manifest.Save(root); err != nil {
		return err
	}

	slog.Info(translator.T("kataStarted"), "kata", kataName)
	return nil

}
