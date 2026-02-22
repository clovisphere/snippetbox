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

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.index))
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(app.show))
	mux.Handle("GET /snippets/create", dynamic.ThenFunc(app.create))
	mux.Handle("POST /snippets/create", dynamic.ThenFunc(app.createPost))

	// standard defines a middleware chain for all application routes.
	// Handlers are executed in order: recoverPanic -> logRequest -> commonHeaders.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
