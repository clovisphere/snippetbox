package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes sets up the application's HTTP routes and returns a ServeMux wrapped
// with standard middleware. It serves static files, defines dynamic routes
// with session management, and enforces authentication for protected routes.
func (app *application) routes() http.Handler {

	// Create a new ServeMux for routing
	mux := http.NewServeMux()

	// Serve static files from ./ui/static/ at /static/
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// dynamic defines a middleware chain for routes that require session state.
	// It automatically loads and saves session data for the current request.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Public routes (no authentication required)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.index))
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(app.show))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userAuthenticate))

	// Protected routes (require authentication)
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippets/create", protected.ThenFunc(app.create))
	mux.Handle("POST /snippets/create", protected.ThenFunc(app.createPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// standard defines a middleware chain for all application routes.
	// Handlers are executed in order: recoverPanic -> logRequest -> commonHeaders.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Wrap all routes with the standard middleware chain and return
	return standard.Then(mux)
}
