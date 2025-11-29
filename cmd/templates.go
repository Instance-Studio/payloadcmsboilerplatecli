package cmd

import "embed"

//go:embed templates/*
var TemplatesFS embed.FS
