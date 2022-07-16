package static

import "embed"

//go:embed testdata/*
var TestFiles embed.FS
