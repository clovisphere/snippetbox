package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/clovisphere/snippetbox/internal/models"
	"github.com/clovisphere/snippetbox/internal/validator"
)

// snippetCreateForm holds the form data and validation errors for a new snippet.
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

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

// create renders the page for creating a new resource.
func (app *application) create(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.html", data)
}

// createPost processes the form submission for creating a new resource.
func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form data against specific business rules.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.storage.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
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
