package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/nadiannis/stn/ui"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	files := []string{
		"html/base.tmpl.html",
		"html/partials/header.tmpl.html",
		"html/partials/footer.tmpl.html",
		"html/pages/home.tmpl.html",
	}
	ts, err := template.ParseFS(ui.Files, files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) linkList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("View all links"))
}

func (app *application) linkCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("View create short link page"))
	case http.MethodPost:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Create short link"))
	default:
		allowedMethods := fmt.Sprintf("%s, %s", http.MethodGet, http.MethodPost)
		w.Header().Set("Allow", allowedMethods)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
