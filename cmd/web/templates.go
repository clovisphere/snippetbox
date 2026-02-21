package main

import (
	"html/template"
	"path/filepath"

	"github.com/clovisphere/snippetbox/internal/models"
)

// templateData holds dynamic data passed to HTML templates.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

// newTemplateCache parses all page templates along with base layout
// and partials, returning a map of template name to compiled template.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
