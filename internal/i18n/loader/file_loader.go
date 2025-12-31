package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
// basePath is the path where the 'i18n' directory is expected (e.g., project root).
func NewFileLoader(basePath string) *FileLoader {
	return &FileLoader{
		basePath: basePath,
	}
}

// Load loads messages for a given language tag from TOML files in the filesystem.
// It expects files in the format basePath/i18n/{lang}.toml (e.g., /project/root/i18n/en.toml).
func (l *FileLoader) Load(tag language.Tag) ([]i18n.Message, error) {
	var messages []i18n.Message

	langCode := strings.SplitN(tag.String(), "-", 2)[0] // e.g., "en-US" -> "en"
	fileName := fmt.Sprintf("%s.toml", langCode)
	fullPath := filepath.Join(l.basePath, "i18n", fileName)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// File not found is not an error, it just means no custom translations exist.
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("checking file status for %s: %w", fullPath, err)
	}

	var rawMessages map[string]string
	if _, err := toml.DecodeFile(fullPath, &rawMessages); err != nil {
		return nil, fmt.Errorf("unmarshaling TOML for %s: %w", fullPath, err)
	}

	for key, value := range rawMessages {
		messages = append(messages, i18n.Message{Key: key, Value: value})
	}

	return messages, nil
}
