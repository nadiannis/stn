package main

import (
	"net/http"
)

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, http.StatusOK, "home.tmpl.html", nil)
}

func (app *application) linkListView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "link-list.tmpl.html", nil)
}

func (app *application) linkCreateView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "link-create.tmpl.html", nil)
}

func (app *application) signupView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "signup.tmpl.html", nil)
}

func (app *application) loginView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "login.tmpl.html", nil)
}
