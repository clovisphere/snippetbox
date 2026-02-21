package main

import (
	"html/template"
	"net/http"
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
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Save a new snippet..."))
}

func (app *application) list(w http.ResponseWriter, r *http.Request) {
	// TODO: show all snippets
}

func (app *application) show(w http.ResponseWriter, r *http.Request) {
	// TODO: display snippet by ID
}
