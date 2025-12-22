package kata

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rjfonseca/kata/internal/assets"
)

// Repository provides access to the kata catalog.
type Repository struct {
	root string
}

// NewRepository creates a new catalog repository.
func NewRepository(root string) *Repository {
	return &Repository{root: root}
}

// ListKatas returns a sorted list of available kata names.
// It prioritizes the local `katas/catalog` directory over the embedded catalog.
func (r *Repository) ListKatas() ([]string, error) {
	localCatalogDir := filepath.Join(r.root, "katas", "catalog")

	var catalogFS fs.FS
	var basePath string

	if _, err := os.Stat(localCatalogDir); err == nil {
		// Local catalog exists
		catalogFS = os.DirFS(localCatalogDir)
		basePath = "."
	} else {
		// Fallback to embedded catalog
		catalogFS = assets.FS
		basePath = "catalog"
	}

	entries, err := fs.ReadDir(catalogFS, basePath)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var katas []string
	for _, e := range entries {
		if e.IsDir() {
			name := e.Name()
			if strings.HasPrefix(name, "_") {
				name = strings.TrimPrefix(name, "_")
			}
			if !seen[name] {
				katas = append(katas, name)
				seen[name] = true
			}
		}
	}

	sort.Strings(katas)
	return katas, nil
}
