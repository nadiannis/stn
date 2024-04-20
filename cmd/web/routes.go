package main

import (
	"net/http"

	"github.com/nadiannis/stn/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", fileServer)

	mux.Handle("GET /", app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.homeView))))
	mux.Handle("GET /{backhalf}", app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.linkRedirect))))
	mux.Handle("GET /users/signup", app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.signupView))))
	mux.Handle("POST /users/signup", app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.signup))))
	mux.Handle("GET /users/login", app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.loginView))))
	mux.Handle("POST /users/login", app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.login))))
	mux.Handle("POST /links/create", app.sessionManager.LoadAndSave(app.authenticate(http.HandlerFunc(app.linkCreate))))

	mux.Handle("POST /users/logout", app.sessionManager.LoadAndSave(app.authenticate(app.protectRoute(http.HandlerFunc(app.logout)))))
	mux.Handle("GET /links/list", app.sessionManager.LoadAndSave(app.authenticate(app.protectRoute(http.HandlerFunc(app.linkListView)))))
	mux.Handle("GET /links/create", app.sessionManager.LoadAndSave(app.authenticate(app.protectRoute(http.HandlerFunc(app.linkCreateView)))))
	mux.Handle("GET /links/{id}", app.sessionManager.LoadAndSave(app.authenticate(app.protectRoute(http.HandlerFunc(app.linkDetailView)))))
	mux.Handle("GET /links/{id}/edit", app.sessionManager.LoadAndSave(app.authenticate(app.protectRoute(http.HandlerFunc(app.linkEditView)))))
	mux.Handle("PUT /links/{id}/edit", app.sessionManager.LoadAndSave(app.authenticate(app.protectRoute(http.HandlerFunc(app.linkEdit)))))

	return app.recoverPanic(app.requestLogger(secureHeaders(mux)))
}
