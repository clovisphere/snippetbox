package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
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

// render looks up the template by name, executes it into a buffer, and
// writes the result to the http.ResponseWriter with the provided status code.
// If the template does not exist or execution fails, a server error is logged
// and a 500 Internal Server Error is sent.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Look up the requested template in the cache
	ts, ok := app.templateCache[page]
	if !ok {
		// Template not found; log and respond with a server error
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Execute the template into a temporary buffer first to catch errors
	buf := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buf, "base", data); err != nil {
		// Template execution failed; log and respond with a server error
		app.serverError(w, r, err)
		return
	}

	// Write the provided HTTP status code
	w.WriteHeader(status)

	// Write the buffered template content to the response
	buf.WriteTo(w)
}

// newTemplateData initializes a templateData struct and populates it with
// common dynamic data: the CSRF token, the current year, any flash message
// in the session, and the user's authentication status.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CSRFToken:       nosurf.Token(r),
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

// decodePostForm parses the request and decodes the post form data into a destination struct.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := app.formDecoder.Decode(dst, r.PostForm); err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}

// isAuthenticated returns true if the current request is from an authenticated user,
// otherwise it returns false.
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
