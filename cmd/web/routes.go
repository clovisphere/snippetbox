package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes sets up the application's HTTP routes and returns a ServeMux.
// It serves static files from ./ui/static/ and maps URL patterns to
// handler methods like index, list, show, and create.
func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// dynamic defines a middleware chain for routes that require session state.
	// It automatically loads and saves session data for the current request.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Snippet routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.index))
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(app.show))
	mux.Handle("GET /snippets/create", dynamic.ThenFunc(app.create))
	mux.Handle("POST /snippets/create", dynamic.ThenFunc(app.createPost))

	// User authentication and registration routes
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userAuthenticate))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// standard defines a middleware chain for all application routes.
	// Handlers are executed in order: recoverPanic -> logRequest -> commonHeaders.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
