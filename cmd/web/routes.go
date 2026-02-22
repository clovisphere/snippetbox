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

	mux.HandleFunc("GET /{$}", app.index)
	mux.HandleFunc("GET /snippets/view/{id}", app.show)
	mux.HandleFunc("GET /snippets/create", app.create)
	mux.HandleFunc("POST /snippets/create", app.createPost)

	// standard defines a middleware chain for all application routes.
	// Handlers are executed in order: recoverPanic -> logRequest -> commonHeaders.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
