package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/clovisphere/snippetbox/internal/models"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/pages/base.html",
		"./ui/html/pages/partials/nav.html",
		"./ui/html/pages/home.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	// TODO: create a snippet
}

func (app *application) list(w http.ResponseWriter, r *http.Request) {
	// TODO: show all snippets
}

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

	fmt.Fprintf(w, "%+v", snippet)
}
