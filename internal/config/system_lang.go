package config

import (
	"os"
	"strings"
)

func SystemLang() string {
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if val := os.Getenv(key); val != "" {
			return normalizeLang(val)
		}
	}
	return ""
}

func normalizeLang(lang string) string {
	lang, _, _ = strings.Cut(lang, ".")
	lang = strings.ReplaceAll(lang, "_", "-")
	return lang
}
