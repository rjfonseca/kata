package config

import "os"

type Resolver struct {
	cfg *Config
}

// Resolver resolves configuration values using precedence rules
// (flags > env > config > system > fallback).
func NewResolver(cfg *Config) *Resolver {
	return &Resolver{cfg: cfg}
}

func (r *Resolver) ResolveLang(flagLang string) string {
	if flagLang != "" {
		return flagLang
	}

	if env := os.Getenv("KATA_LANG"); env != "" {
		return normalizeLang(env)
	}

	if r.cfg.UI.Lang != "" {
		return r.cfg.UI.Lang
	}

	if sys := SystemLang(); sys != "" {
		return sys
	}

	return "en-US"
}
