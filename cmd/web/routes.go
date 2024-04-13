package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/links/list", app.linkList)
	mux.HandleFunc("/links/create", app.linkCreate)

	return app.requestLogger(mux)
}
