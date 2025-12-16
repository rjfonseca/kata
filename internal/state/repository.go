package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Repository manages persistence of the kata state on disk.
type Repository struct {
	path string
}

// NewRepository creates a new state repository rooted at .kata/current_state.json.
func NewRepository(root string) *Repository {
	return &Repository{
		path: filepath.Join(root, ".kata", "current_state.json"),
	}
}

// Exists returns true if a kata state already exists.
func (r *Repository) Exists() bool {
	_, err := os.Stat(r.path)
	return err == nil
}

// Load loads the kata state from disk.
func (r *Repository) Load() (*State, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

// Save persists the kata state to disk, creating directories if needed.
func (r *Repository) Save(state *State) error {
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, data, 0o644)
}

// Delete removes the kata state from disk.
func (r *Repository) Delete() error {
	if !r.Exists() {
		return errors.New("state does not exist")
	}
	return os.Remove(r.path)
}

// FindProjectRoot searches for the .kata directory starting from startDir and moving up.
func FindProjectRoot(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		if info, err := os.Stat(filepath.Join(dir, ".kata")); err == nil && info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("kata state not found (are you in a kata project?)")
		}
		dir = parent
	}
}
