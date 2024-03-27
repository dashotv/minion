package static

import "embed"

//go:embed *.html all:assets/*
var FS embed.FS
