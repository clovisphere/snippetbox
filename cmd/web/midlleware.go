package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

// commonHeaders is middleware that adds security-focused headers (CSP, Referrer,
// and Frame options) to all outgoing responses.
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		h.Set("Referrer-Policy", "origin-when-cross-origin")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "deny")
		// Disable the outdated XSS filter in favor of the CSP policy.
		h.Set("X-XSS-Protection", "0")
		h.Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// logRequest is middleware that logs metadata about every HTTP request,
// including the remote IP, protocol, HTTP method, and requested URI.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info(
			"Received request",
			slog.String("ip", ip),
			slog.String("proto", proto),
			slog.String("method", method),
			slog.String("uri", uri),
			slog.String("user-agent", r.UserAgent()),
		)

		next.ServeHTTP(w, r)
	})
}

// recoverPanic is middleware that recovers from any panics during request handling,
// logs the error, and sends a 500 Internal Server Error response.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				// Signal that the connection should be closed after this response
				w.Header().Set("Connection", "close")

				app.serverError(w, r, fmt.Errorf("%v", p))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// requireAuthentication is middleware that ensures a user is authenticated
// before accessing the next handler. If the user is not authenticated, it
// redirects them to the login page. It also sets headers to prevent caching
// of sensitive content.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect to login if the user is not authenticated
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Prevent caching of authenticated pages
		w.Header().Add("Cache-Control", "no-store")

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
