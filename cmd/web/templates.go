package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/clovisphere/snippetbox/internal/models"
	"github.com/clovisphere/snippetbox/ui"
)

// templateData acts as a container for any dynamic data that we want
// to pass to our HTML templates.
type templateData struct {
	CSRFToken       string           // A token to prevent CSRF attacks on POST forms.
	CurrentYear     int              // The current year for the footer copyright notice.
	Flash           string           // A one-time message for the user (e.g., "Login successful").
	Form            any              // Form data and validation errors for re-population.
	IsAuthenticated bool             // A boolean flag to toggle UI elements for logged-in users.
	Snippet         models.Snippet   // A single snippet record.
	Snippets        []models.Snippet // A slice of multiple snippet records.
}

// functions defines custom template functions available in HTML templates.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// humanDate returns a formatted string representation of a time.Time value
// in the form "02 Jan 2006 at 15:04".
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// newTemplateCache creates a map of compiled templates, allowing for fast
// lookups during request handling. It uses an embedded file system to
// ensure templates are bundled within the application binary.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use fs.Glob to find all files in the pages directory within the embedded FS.
	// Note: paths in embed.FS do not use the './' prefix.
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Define the template patterns to parse. The order matters:
		// base.html defines the layout, partials provide reusable components,
		// and the page template fills in the specific content.
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		// ParseFS replaces ParseFiles/ParseGlob for embedded files.
		// It's an atomic way to compile the entire template set for a page.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
