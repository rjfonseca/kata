package assets

import "embed"

// FS contains embedded assets used to bootstrap kata repositories.
//
//go:embed catalog/**/* runners/**/* i18n/* kata.toml
var FS embed.FS
