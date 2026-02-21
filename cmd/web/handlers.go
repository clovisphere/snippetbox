package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/clovisphere/snippetbox/internal/models"
)

// index fetches the latest snippets from storage and renders the home page.
func (app *application) index(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.storage.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Populate template data with the latest snippets and render the home page.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(
		w,
		r,
		http.StatusOK,
		"home.html",
		data,
	)
}

// create is a placeholder handler for creating a new snippet.
func (app *application) create(w http.ResponseWriter, r *http.Request) {
	// TODO: create a snippet
}

// list is a placeholder handler for listing all snippets.
func (app *application) list(w http.ResponseWriter, r *http.Request) {
	// TODO: show all snippets
}

// show fetches a snippet by ID from storage and renders the view page.
// If the ID is invalid or the snippet does not exist, it returns a 404.
func (app *application) show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.storage.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Populate template data with the snippet and render the view page.
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(
		w,
		r,
		http.StatusOK,
		"view.html",
		data,
	)
}
