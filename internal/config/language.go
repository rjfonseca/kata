package config

import (
	"log"

	"github.com/Xuanwo/go-locale"
	"golang.org/x/text/language"
)

// ResolveLanguage determina a tag de idioma a ser usada com base em uma hierarquia de fontes:
// 1. A flag da linha de comando (`langFlag`).
// 2. O locale do sistema operacional.
// 3. Inglês como padrão final.
func ResolveLanguage(langFlag string) language.Tag {
	// 1. Prioridade máxima: flag da linha de comando
	if langFlag != "" {
		tag, err := language.Parse(langFlag)
		if err != nil {
			log.Printf("Warning: Could not parse language flag '%s', using default. Error: %v\n", langFlag, err)
			return language.English
		}
		return tag
	}

	// 2. Prioridade média: locale do sistema
	tag, err := locale.Detect()
	if err == nil {
		return tag
	}

	// Loga o erro de detecção, mas continua para o fallback
	log.Printf("Warning: Could not detect OS locale, using default. Error: %v\n", err)

	// 3. Padrão: Inglês
	return language.English
}
