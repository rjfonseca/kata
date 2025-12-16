package scaffold

import "testing"

func TestManifestLifecycle(t *testing.T) {
	tmp := t.TempDir()

	m, err := LoadManifest(tmp)
	if err != nil {
		t.Fatal(err)
	}

	m.Files["foo.txt"] = FileEntry{
		Checksum: "sha256:abc",
		Source:   "scaffold",
	}

	if err := m.Save(tmp); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadManifest(tmp)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := loaded.Files["foo.txt"]; !ok {
		t.Fatalf("expected file entry to persist")
	}
}
