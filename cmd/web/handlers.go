package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/clovisphere/snippetbox/internal/models"
	"github.com/clovisphere/snippetbox/internal/validator"
)

// userLoginForm holds the form data and validation errors for the user login flow.
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// userSignupForm holds the form data and validation errors for user registration.
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// snippetCreateForm holds the form data and validation errors for a new snippet.
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// index fetches the latest snippets from storage and renders the home page.
func (app *application) index(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
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

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

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

	snippet, err := app.snippets.Get(id)
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

// userSignup handles the GET request to display the user registration form.
// It initializes a new userSignupForm and renders the signup template.
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form data against specific business rules.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRx), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field mujst be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	if err := app.users.Insert(form.Name, form.Email, form.Password); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// userLogin handles the GET request to display the login form.
// It initializes a new userLoginForm and renders the login template.
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) userAuthenticate(w http.ResponseWriter, r *http.Request) {
	// Verifies credentials and creates a session
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Destroys the session
}
