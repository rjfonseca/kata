package i18n

import (
	"fmt"
	"log/slog"

	"golang.org/x/text/language"
)

// Translator define o contrato para tradução de mensagens.
// É a única interface que o resto da aplicação deve conhecer.
type Translator interface {
	// T traduz uma mensagem para a chave fornecida, formatando-a com os argumentos.
	// Retorna a string traduzida ou uma string de fallback caso a chave não seja encontrada.
	T(key string, args ...any) string
}

// DefaultTranslator é a implementação concreta da interface Translator.
type DefaultTranslator struct {
	// messages armazena as traduções carregadas, key -> translated_string
	messages map[string]string
	// fallbackMessages armazena as traduções do idioma de fallback (geralmente inglês)
	fallbackMessages map[string]string
	// currentLang indica o idioma atual do tradutor
	currentLang language.Tag
}

// NewDefaultTranslator cria e inicializa um novo DefaultTranslator.
// Ele recebe um MessageLoader que será usado para carregar as mensagens.
// O idioma preferencial é detectado a partir do ambiente ou passado como padrão.
func NewDefaultTranslator(loader MessageLoader, preferredLangs ...language.Tag) (*DefaultTranslator, error) {
	// Dynamically get the list of supported languages from the embedded map
	var supportedTags []language.Tag
	for langStr := range EmbeddedLocales {
		supportedTags = append(supportedTags, language.Make(langStr))
	}
	if len(supportedTags) == 0 {
		supportedTags = append(supportedTags, language.English)
	}

	// Create a matcher with the supported languages
	matcher := language.NewMatcher(supportedTags)

	// Determine the desired language tag
	desiredTag := language.English
	if len(preferredLangs) > 0 {
		desiredTag = preferredLangs[0]
	}

	// Use the matcher to find the best supported language tag
	_, index, _ := matcher.Match(desiredTag)
	langTag := supportedTags[index]

	// Load fallback messages (English)
	fallbackMsgs, err := loader.Load(language.English)
	if err != nil {
		return nil, fmt.Errorf("failed to load fallback (English) messages: %w", err)
	}
	fallbackMap := make(map[string]string)
	for _, msg := range fallbackMsgs {
		fallbackMap[msg.Key] = msg.Value
	}

	// Load messages for the matched language
	currentMsgs, err := loader.Load(langTag)
	if err != nil {
		slog.Warn("Could not load messages, using fallback", "lang", langTag, "error", err)
	}

	currentMap := make(map[string]string)
	for k, v := range fallbackMap {
		currentMap[k] = v
	}
	for _, msg := range currentMsgs {
		currentMap[msg.Key] = msg.Value
	}

	return &DefaultTranslator{
		messages:         currentMap,
		fallbackMessages: fallbackMap,
		currentLang:      langTag,
	}, nil
}

// T traduz uma mensagem para a chave fornecida, formatando-a com os argumentos.
// Retorna a string traduzida ou a string de fallback se não encontrada.
func (t *DefaultTranslator) T(key string, args ...any) string {
	if msg, ok := t.messages[key]; ok {
		return fmt.Sprintf(msg, args...)
	}
	// Fallback para o idioma padrão
	if msg, ok := t.fallbackMessages[key]; ok {
		return fmt.Sprintf(msg, args...)
	}
	// Se nem no fallback encontrar, retorna a chave entre parênteses
	return fmt.Sprintf("[MISSING_TRANSLATION:%s]", key)
}
