package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/rjfonseca/kata/internal/config"
	"github.com/rjfonseca/kata/internal/fsutil"
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/scaffold"
	"github.com/rjfonseca/kata/internal/state"
)

type ResetFlags struct {
	Hard   bool
	ToStep int // -1 means use default (current step or hard)
}

func Reset(root string, stateRepo *state.Repository, f ResetFlags, translator i18n.Translator) error {
	st, err := stateRepo.Load()
	if err != nil {
		return err
	}

	wsCfg, err := config.LoadWorkspaceConfig(root)
	if err != nil {
		return err
	}
	protected := wsCfg.GetProtectedPaths()

	// 1. Safety Backup
	backupDir := filepath.Join(root, ".kata", "backups")
	_ = os.MkdirAll(backupDir, 0o755)
	backupPath := filepath.Join(backupDir, fmt.Sprintf("safety-%d.tar.gz", time.Now().Unix()))

	// When backing up, we ignore the .kata directory to avoid recursive backups
	if err := fsutil.CreateArchive(root, backupPath, []string{".kata"}); err != nil {
		slog.Warn("Failed to create safety backup", "error", err)
	} else {
		slog.Info(translator.T("reset.log_safety_backup"), "path", backupPath)
	}

	// 2. Determine target step
	targetStepIndex := st.CurrentStepIndex
	if f.Hard {
		targetStepIndex = -1
	} else if f.ToStep != -10 {
		if f.ToStep >= -1 && f.ToStep < len(st.Steps) {
			targetStepIndex = f.ToStep
		} else {
			return errors.New(translator.T("reset.error_invalid_step"))
		}
	}

	// 3. Find Checkpoint
	// The checkpoint to restore is the one that was created when the PREVIOUS step was finished.
	// If target is step N, we need the checkpoint from step N-1.
	checkpointName := "scaffold"
	if targetStepIndex > 0 {
		prevIndex := targetStepIndex - 1
		checkpointName = fmt.Sprintf("step-%02d-%s", prevIndex, st.Steps[prevIndex])
	} else if targetStepIndex == 0 {
		checkpointName = "scaffold"
	}

	checkpointPath := filepath.Join(root, ".kata", "checkpoints", checkpointName+".tar.gz")
	if _, err := os.Stat(checkpointPath); os.IsNotExist(err) {
		return fmt.Errorf("%s: %s", translator.T("reset.error_checkpoint_not_found"), checkpointName)
	}

	slog.Info(translator.T("reset.log_restoring"), "checkpoint", checkpointName)

	// 4. Clean Workspace
	if err := fsutil.CleanDir(root, protected); err != nil {
		return err
	}

	// 5. Restore Checkpoint
	if err := fsutil.ExtractArchive(checkpointPath, root); err != nil {
		return err
	}

	// 6. If not a hard reset to scaffold, re-apply the target step's challenge
	if targetStepIndex >= 0 {
		stepName := st.Steps[targetStepIndex]
		stepDir := filepath.Join(".kata", "catalog", st.KataName, "steps", stepName)

		manifest, err := scaffold.LoadManifest(root)
		if err == nil {
			copier := &scaffold.Copier{
				Root:       root,
				Manifest:   manifest,
				Translator: translator,
			}
			_ = copier.Apply(os.DirFS(stepDir), "step:"+stepName)
			_ = manifest.Save(root)
		}
	} else {
		// If hard reset, we might want to re-apply the kata scaffold to be sure
		scaffoldDir := filepath.Join(".kata", "catalog", st.KataName, "scaffold")
		if info, err := os.Stat(scaffoldDir); err == nil && info.IsDir() {
			manifest, _ := scaffold.LoadManifest(root)
			copier := &scaffold.Copier{
				Root:       root,
				Manifest:   manifest,
				Translator: translator,
			}
			_ = copier.Apply(os.DirFS(scaffoldDir), "scaffold")
			_ = manifest.Save(root)
		}
	}

	// 7. Update State
	st.CurrentStepIndex = targetStepIndex
	st.TestPassing = false
	st.KataFinished = false
	if err := stateRepo.Save(st); err != nil {
		return err
	}

	currentStepName := "scaffold"
	if st.CurrentStepIndex >= 0 {
		currentStepName = st.Steps[st.CurrentStepIndex]
	}
	slog.Info(translator.T("reset.log_reset_complete"), "step", currentStepName)

	return nil
}
