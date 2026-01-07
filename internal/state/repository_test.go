package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRepository_SaveAndLoad(t *testing.T) {
	tmp := t.TempDir()

	repo := NewRepository(tmp)

	original := &State{
		KataName:         "hello-world",
		StartedAt:        time.Now(),
		Steps:            []string{"01-init", "02-next"},
		CurrentStepIndex: 0,
	}

	if err := repo.Save(original); err != nil {
		t.Fatalf("failed to save state: %v", err)
	}

	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("failed to load state: %v", err)
	}

	if loaded.KataName != original.KataName {
		t.Fatalf("expected kata %q, got %q", original.KataName, loaded.KataName)
	}

	if len(loaded.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(loaded.Steps))
	}

	expectedPath := filepath.Join(tmp, ".kata", "current_state.json")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Fatalf("state file not found at expected path")
	}
}
