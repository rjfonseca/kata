package loader

import (
	"strings"

	"golang.org/x/text/language"

	"github.com/rjfonseca/kata/internal/i18n"
)

// EmbeddedLoader implements the i18n.MessageLoader interface
// by reading messages from a pre-generated Go map.
type EmbeddedLoader struct{}

// NewEmbeddedLoader creates a new EmbeddedLoader instance.
func NewEmbeddedLoader() *EmbeddedLoader {
	return &EmbeddedLoader{}
}

// Load loads messages for a given language tag from the generated map.
func (l *EmbeddedLoader) Load(tag language.Tag) ([]i18n.Message, error) {
	langCode := strings.SplitN(tag.String(), "-", 2)[0] // e.g., "en-US" -> "en"

	langMessages, ok := i18n.EmbeddedLocales[langCode]
	if !ok {
		// No translations for this language, not an error.
		return nil, nil
	}

	messages := make([]i18n.Message, 0, len(langMessages))
	for key, value := range langMessages {
		messages = append(messages, i18n.Message{Key: key, Value: value})
	}

	return messages, nil
}
