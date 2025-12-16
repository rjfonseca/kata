package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var ErrNoWorkspace = errors.New("workspace configuration (kata.toml) not found")

// FindWorkspaceRoot searches for kata.toml starting from startDir and moving up.
func FindWorkspaceRoot(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "kata.toml")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNoWorkspace
		}
		dir = parent
	}
}

// WorkspaceConfig represents the configuration stored in `kata.toml` at the project root.
type WorkspaceConfig struct {
	Workspace struct {
		ProtectedPaths []string `toml:"protected_paths"`
	} `toml:"workspace"`
}

// LoadWorkspaceConfig attempts to load `kata.toml` from the given root directory.
// If the file doesn't exist, it returns a default configuration.
func LoadWorkspaceConfig(root string) (*WorkspaceConfig, error) {
	path := filepath.Join(root, "kata.toml")
	cfg := &WorkspaceConfig{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	} else if err != nil {
		return nil, err
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetProtectedPaths returns the combined list of system-protected paths and user-defined ones.
func (c *WorkspaceConfig) GetProtectedPaths() []string {
	// Paths that are always protected by the system
	systemProtected := []string{
		".git",
		".kata",
		"katas",
		"kata.toml",
		".task",
		".kata_backups",
	}

	return append(systemProtected, c.Workspace.ProtectedPaths...)
}
