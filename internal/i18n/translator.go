package i18n

import (
	"fmt"

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
	// 1. Detectar o idioma preferencial
	// Por enquanto, vamos assumir 'en' como padrão se nada for fornecido ou detectado.
	// A detecção real virá de variáveis de ambiente.
	var langTag language.Tag
	if len(preferredLangs) > 0 {
		langTag = preferredLangs[0]
	} else {
		// Placeholder para detecção real do sistema operacional.
		// Por enquanto, usaremos "en" como padrão.
		langTag = language.English
	}

	// 2. Carregar mensagens de fallback (inglês)
	fallbackMsgs, err := loader.Load(language.English)
	if err != nil {
		return nil, fmt.Errorf("failed to load fallback (English) messages: %w", err)
	}
	fallbackMap := make(map[string]string)
	for _, msg := range fallbackMsgs {
		fallbackMap[msg.Key] = msg.Value
	}

	// 3. Carregar mensagens para o idioma atual
	currentMsgs, err := loader.Load(langTag)
	if err != nil {
		// Não é um erro crítico se não houver tradução para o idioma atual,
		// apenas usaremos o fallback.
		fmt.Printf("Warning: Could not load messages for language %s, using fallback. Error: %v\n", langTag, err)
	}
	currentMap := make(map[string]string)
	// Primeiro preenche com o fallback, depois sobrescreve com o idioma atual
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
	},
	nil
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
