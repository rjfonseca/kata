package scaffold

import (
	"os"
	"testing"
	"testing/fstest"
)

func TestCopierWritesAndTracksFile(t *testing.T) {
	fs := fstest.MapFS{
		"foo.txt": {Data: []byte("hello")},
	}

	tmp := t.TempDir()
	m, _ := LoadManifest(tmp)

	c := &Copier{
		Root:     tmp,
		Manifest: m,
	}

	if err := c.Apply(fs, "scaffold"); err != nil {
		t.Fatal(err)
	}

	path := tmp + "/foo.txt"
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created")
	}

	if _, ok := m.Files[path]; !ok {
		t.Fatalf("file not tracked in manifest")
	}
}
