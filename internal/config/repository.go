package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Repository handles loading and saving the global kata configuration
// from the user's config directory.
type Repository struct {
	path string
}

// NewRepository creates a new configuration repository using the
// OS-specific user config directory.
func NewRepository() (*Repository, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "kata", "config.toml")
	return &Repository{path: path}, nil
}

func (r *Repository) Load() (*Config, error) {
	if _, err := os.Stat(r.path); errors.Is(err, os.ErrNotExist) {
		cfg := DefaultConfig()
		if err := r.Save(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	var cfg Config
	if _, err := toml.DecodeFile(r.path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (r *Repository) Save(cfg *Config) error {
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f, err := os.Create(r.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}
