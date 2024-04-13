package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"text/template"

	"github.com/nadiannis/stn/ui"
)

var templateFiles = map[string]string{
	"home":        "html/pages/home.tmpl.html",
	"link-list":   "html/pages/link-list.tmpl.html",
	"link-create": "html/pages/link-create.tmpl.html",
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	files := []string{
		"html/base.tmpl.html",
		"html/partials/header.tmpl.html",
		"html/partials/footer.tmpl.html",
		templateFiles[page],
	}
	ts, err := template.ParseFS(ui.Files, files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}
