package config

import (
	"os"
	"testing"
)

func TestResolveLang_FromFlag(t *testing.T) {
	cfg := DefaultConfig()
	r := NewResolver(cfg)

	lang := r.ResolveLang("pt-BR")
	if lang != "pt-BR" {
		t.Fatalf("expected pt-BR, got %s", lang)
	}
}

func TestResolveLang_FromEnv(t *testing.T) {
	t.Setenv("KATA_LANG", "en_US.UTF-8")

	cfg := DefaultConfig()
	r := NewResolver(cfg)

	lang := r.ResolveLang("")
	if lang != "en-US" {
		t.Fatalf("expected en-US, got %s", lang)
	}
}

func TestResolveLang_Fallback(t *testing.T) {
	os.Unsetenv("KATA_LANG")

	cfg := DefaultConfig()
	r := NewResolver(cfg)

	lang := r.ResolveLang("")
	if lang == "" {
		t.Fatalf("expected fallback lang, got empty")
	}
}
