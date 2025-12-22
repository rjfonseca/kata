package catalog

import (
	"errors"
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/rjfonseca/kata/internal/assets"
)

// Definition represents a kata definition loaded from the catalog.
type Definition struct {
	Name  string
	Steps []string
}

// Load loads a kata definition from the embedded catalog.
func Load(name string) (*Definition, error) {
	basePath := filepath.Join("catalog", name)

	entries, err := fs.ReadDir(assets.FS, filepath.Join(basePath, "steps"))
	if err != nil {
		return nil, errors.New("kata not found in catalog")
	}

	var steps []string
	for _, e := range entries {
		if e.IsDir() {
			steps = append(steps, e.Name())
		}
	}

	sort.Strings(steps)

	return &Definition{
		Name:  name,
		Steps: steps,
	}, nil
}
