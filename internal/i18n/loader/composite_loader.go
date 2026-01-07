package loader

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/rjfonseca/kata/internal/i18n"
)

// CompositeLoader implements the i18n.MessageLoader interface
// by combining multiple loaders. It applies an override strategy
// where messages from later loaders in the slice take precedence.
type CompositeLoader struct {
	loaders []i18n.MessageLoader
}

// NewCompositeLoader creates a new CompositeLoader instance.
// Loaders are processed in order, with later loaders overriding
// messages from earlier ones for the same key.
func NewCompositeLoader(loaders ...i18n.MessageLoader) *CompositeLoader {
	return &CompositeLoader{
		loaders: loaders,
	}
}

// Load loads messages from all configured loaders, applying overrides.
func (l *CompositeLoader) Load(tag language.Tag) ([]i18n.Message, error) {
	allMessages := make(map[string]string)

	for i, loader := range l.loaders {
		currentLoaderMessages, err := loader.Load(tag)
		if err != nil {
			// If a loader fails, we might want to log it and continue,
			// or return an error depending on desired strictness.
			// For now, let's return an error as failing to load any part
			// of the translation chain might indicate a bigger issue.
			return nil, fmt.Errorf("failed to load messages from loader %d: %w", i, err)
		}

		for _, msg := range currentLoaderMessages {
			allMessages[msg.Key] = msg.Value // Overwrites previous messages for the same key
		}
	}

	var result []i18n.Message
	for key, value := range allMessages {
		result = append(result, i18n.Message{Key: key, Value: value})
	}

	return result, nil
}
