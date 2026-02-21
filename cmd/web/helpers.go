package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// serverError logs an internal server error along with the request method,
// URI, and stack trace, then sends a 500 Internal Server Error response.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(
		err.Error(),
		slog.String("method", method),
		slog.String("uri", uri),
		slog.String("trace", trace),
	)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a response with the given HTTP status code and its
// associated status text. Useful for 4xx client errors.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// render looks up the template by name, writes the provided status code,
// and executes the template with the provided data. If the template does
// not exist or execution fails, it logs a server error.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	if err := ts.ExecuteTemplate(w, "base", data); err != nil {
		app.serverError(w, r, err)
	}
}
