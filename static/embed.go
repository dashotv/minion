package static

import "embed"

//go:embed index.html all:assets/*
var FS embed.FS
