package main

import "embed"

//go:embed templates/index.html
var indexHTML string

//go:embed static/*
var staticFiles embed.FS
