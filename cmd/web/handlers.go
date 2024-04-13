package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, http.StatusOK, "home", nil)
}

func (app *application) linkList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.render(w, http.StatusOK, "link-list", nil)
}

func (app *application) linkCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.render(w, http.StatusOK, "link-create", nil)
	case http.MethodPost:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Create short link"))
	default:
		allowedMethods := fmt.Sprintf("%s, %s", http.MethodGet, http.MethodPost)
		w.Header().Set("Allow", allowedMethods)
		app.clientError(w, http.StatusMethodNotAllowed)
	}
}
