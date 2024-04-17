package main

import (
	"net/http"

	"github.com/nadiannis/stn/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", fileServer)

	mux.Handle("GET /", app.sessionManager.LoadAndSave(http.HandlerFunc(app.homeView)))
	mux.Handle("GET /links/list", app.sessionManager.LoadAndSave(http.HandlerFunc(app.linkListView)))
	mux.Handle("GET /links/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.linkCreateView)))
	mux.Handle("GET /links/{id}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.linkDetailView)))
	mux.Handle("GET /links/{id}/edit", app.sessionManager.LoadAndSave(http.HandlerFunc(app.linkEditView)))
	mux.Handle("GET /users/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.signupView)))
	mux.Handle("POST /users/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.signup)))
	mux.Handle("GET /users/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.loginView)))
	mux.Handle("POST /users/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.login)))

	return app.requestLogger(secureHeaders(mux))
}
