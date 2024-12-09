package v5

import (
	"embed"
	_ "embed"
)

//go:embed testdata/5.yaml
var testData5 string

//go:embed testdata/*
var testFiles embed.FS
