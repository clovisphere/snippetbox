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

// newTemplateCache parses all page templates along with the base layout
// and partials, returning a map of template name to compiled template.
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize an empty map to store compiled templates
	cache := map[string]*template.Template{}

	// Find all page templates in the pages directory
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop over each page template
	for _, page := range pages {
		// Extract the file name (e.g., "home.html") to use as the cache key
		name := filepath.Base(page)

		// Start by parsing the base layout template
		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// Parse all partial templates (e.g., nav, footer) into the same template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Parse the current page template, so it can override blocks in base
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the fully compiled template set to the cache with the page name as key
		cache[name] = ts
	}

	// Return the cache containing all compiled templates
	return cache, nil
}
