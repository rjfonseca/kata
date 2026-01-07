package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.UI.Lang != "" {
		t.Fatalf("expected empty lang, got %q", cfg.UI.Lang)
	}

	if cfg.UI.NonInteractive {
		t.Fatalf("expected non_interactive to be false")
	}
}
