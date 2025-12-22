package fsutil

import (
	"os"
	"path/filepath"
)

// CleanDir removes all files and directories in the given root directory,
// except those specified in the protectedPaths list.
func CleanDir(root string, protectedPaths []string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	protectedMap := make(map[string]bool)
	for _, p := range protectedPaths {
		protectedMap[p] = true
	}

	for _, e := range entries {
		if protectedMap[e.Name()] {
			continue
		}

		path := filepath.Join(root, e.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}
