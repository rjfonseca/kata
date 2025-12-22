package fsutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rjfonseca/kata/internal/assets"
)

func TestCopyDir(t *testing.T) {
	tmp := t.TempDir()

	err := CopyDir(assets.FS, "catalog", filepath.Join(tmp, "catalog"))
	if err != nil {
		t.Fatalf("copy failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmp, "catalog", "hello-world")); err != nil {
		t.Fatalf("expected catalog to be copied")
	}
}
