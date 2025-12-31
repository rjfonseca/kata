// Package i18n provides the internationalization capabilities for the CLI.
//
//go:generate sh -c "go build -o /tmp/i18n-codegen github.com/rjfonseca/kata/tools/i18n-codegen && /tmp/i18n-codegen -in locales -out locales.go"
package i18n
