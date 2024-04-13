package main

import (
	"net/http"

	"github.com/nadiannis/stn/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("/static/", fileServer)

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/links/list", app.linkList)
	mux.HandleFunc("/links/create", app.linkCreate)

	return app.requestLogger(mux)
}
