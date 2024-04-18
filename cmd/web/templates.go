package main

import (
	"io/fs"
	"path/filepath"
	"text/template"

	"github.com/nadiannis/stn/ui"
)

type templateData struct {
	Form              any
	IsAuthenticated   bool
	AuthenticatedUser *AuthenticatedUser
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

		ts, err := template.ParseFS(ui.Files, files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
