package scaffold_test

import (
	"os"
	"path/filepath"
	"testing"

	"testing/fstest"

	"github.com/rjfonseca/kata/internal/scaffold"
	"github.com/rjfonseca/kata/internal/state"
)

func TestNextAppliesStepFiles(t *testing.T) {
	tmp := t.TempDir()

	// Simula .kata/catalog
	stepFS := fstest.MapFS{
		"bar.txt": {Data: []byte("step2")},
	}

	_ = os.MkdirAll(filepath.Join(tmp, ".kata", "catalog", "kata", "steps", "02"), 0o755)

	// Prepara estado
	s := &state.State{
		KataName:         "kata",
		Steps:            []string{"01", "02"},
		CurrentStepIndex: 0,
	}

	stateRepo := state.NewRepository(tmp)
	_ = stateRepo.Save(s)

	manifest, _ := scaffold.LoadManifest(tmp)
	copier := &scaffold.Copier{Root: tmp, Manifest: manifest}

	// Aplica step 02
	if err := copier.Apply(stepFS, "step:02"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(tmp, "bar.txt")); err != nil {
		t.Fatalf("expected step file to be applied")
	}
}
