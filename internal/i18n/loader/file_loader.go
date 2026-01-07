package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"golang.org/x/text/language"

	"github.com/rjfonseca/kata/internal/i18n"
)

// FileLoader implements the i18n.MessageLoader interface
// to load messages from TOML files in the filesystem.
type FileLoader struct {
	basePath string // The root directory to search for i18n files
}

// NewFileLoader creates a new FileLoader instance.
// basePath is the root of the project.
func NewFileLoader(basePath string) *FileLoader {
	return &FileLoader{
		basePath: filepath.Join(basePath, "katas", "i18n"),
	}
}

// Load loads messages for a given language tag from TOML files in the filesystem.
// It expects files in the format basePath/{lang}.toml (e.g., /project/root/i18n/en.toml).
func (l *FileLoader) Load(tag language.Tag) ([]i18n.Message, error) {
	langCode := tag.String()
	fileName := fmt.Sprintf("%s.toml", langCode)
	fullPath := filepath.Join(l.basePath, fileName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, nil // File not found is not an error
	} else if err != nil {
		return nil, fmt.Errorf("checking file status for %s: %w", fullPath, err)
	}

	var domainMessages map[string]map[string]string
	if _, err := toml.DecodeFile(fullPath, &domainMessages); err != nil {
		return nil, fmt.Errorf("unmarshaling TOML for %s: %w", fullPath, err)
	}

	return flatten(domainMessages), nil
}

func flatten(data map[string]map[string]string) []i18n.Message {
	var messages []i18n.Message
	for domain, msgs := range data {
		for key, value := range msgs {
			messages = append(messages, i18n.Message{Key: domain + "." + key, Value: value})
		}
	}
	return messages
}
