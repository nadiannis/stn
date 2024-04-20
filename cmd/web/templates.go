package main

import (
	"io/fs"
	"path/filepath"
	"text/template"
	"time"

	"github.com/nadiannis/stn/internal/models"
	"github.com/nadiannis/stn/ui"
)

type templateData struct {
	Form              any
	IsAuthenticated   bool
	AuthenticatedUser *AuthenticatedUser
	Flash             string
	Origin            map[string]string
	Link              *models.Link
	Links             []*models.Link
}

func formatDate(t time.Time) string {
	format := "Jan 02, 2006 03:04 PM"
	location, err := time.LoadLocation("Local")
	if err != nil {
		return t.Format(format)
	}
	return t.In(location).Format(format)
}

var functions = template.FuncMap{
	"formatDate": formatDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
