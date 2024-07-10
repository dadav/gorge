package ui

import (
	"embed"
	"net/http"
)

//go:embed all:assets
var assets embed.FS

func HandleAssets() http.Handler {
	return http.FileServer(http.FS(assets))
}
