package main

import (
	"embed"

	"github.com/Gokul-Eswar/Spectre/internal/cli"
	"github.com/Gokul-Eswar/Spectre/internal/server"
)

//go:embed web/*
var WebAssets embed.FS

func main() {
	server.SetAssets(WebAssets)
	cli.Execute()
}
