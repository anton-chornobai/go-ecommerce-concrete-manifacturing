package handlers

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed swagger/*
var swaggerUI embed.FS

func SwaggerDocHandler() http.Handler {
	subFS, err := fs.Sub(swaggerUI, "swagger")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(subFS))
}