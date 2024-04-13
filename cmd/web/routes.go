package main

import (
	"net/http"

	"github.com/nadiannis/stn/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /", app.homeView)
	mux.HandleFunc("GET /links/list", app.linkListView)
	mux.HandleFunc("GET /links/create", app.linkCreateView)
	mux.HandleFunc("GET /links/{id}", app.linkDetailView)
	mux.HandleFunc("GET /users/signup", app.signupView)
	mux.HandleFunc("GET /users/login", app.loginView)

	return app.requestLogger(mux)
}
