package assets

import "embed"

// FS contains embedded assets used to bootstrap kata repositories.
//
//go:embed catalog/**/* runners/**/*
var FS embed.FS
