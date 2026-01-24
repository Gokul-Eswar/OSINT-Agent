package main

import (
	"embed"

	"github.com/spectre/spectre/internal/cli"
	"github.com/spectre/spectre/internal/server"
)

//go:embed web/*
var WebAssets embed.FS

func main() {
	server.SetAssets(WebAssets)
	cli.Execute()
}