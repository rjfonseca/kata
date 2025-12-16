package i18n

import "golang.org/x/text/language"

// Message representa um par chave-valor de tradução.
type Message struct {
	Key   string
	Value string
}

// MessageLoader define o contrato para carregar mensagens de uma fonte específica.
type MessageLoader interface {
	Load(tag language.Tag) (messages []Message, err error)
}
