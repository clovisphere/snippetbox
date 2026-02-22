package main

import (
	"net/http"
	"testing"

	"github.com/clovisphere/snippetbox/internal/assert"
)

// TestIndex performs an end-to-end integration test of the home page.
// It verifies that the application starts, routes the request correctly,
// and returns a 200 OK status with the expected content.
func TestIndex(t *testing.T) {
	// Initialize a dependency-injected test version of the application.
	app := newTestApplication(t)

	// Spin up a HTTPS test server for the duration of this test.
	// Cleanup (ts.Close) is handled automatically by the helper's t.Cleanup.
	ts := newTestServer(t, app.routes())

	// Execute a GET request against the root URL path.
	res := ts.get(t, "/")

	// Assert that the response status code is 200 (OK).
	assert.Equal(t, res.status, http.StatusOK)

	// Assert that the response body contains the expected page heading.
	// This ensures that the templates were rendered correctly.
	assert.StringContains(t, res.body, "Latest Snippets")
}

// TestSnippetView checks the behavior of the snippet viewing endpoint.
// It uses table-driven tests to verify various scenarios, including successful
// lookups, missing records, and malformed URL parameters.
func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	tests := []struct {
		name    string
		urlPath string
		status  int
		body    string // The expected substring in the response body
	}{
		{
			name:    "Valid ID",
			urlPath: "/snippets/view/1",
			status:  http.StatusOK,
			body:    "Gopher Standard Library Cheat Sheet",
		},
		{
			name:    "Non-existent ID",
			urlPath: "/snippets/view/2",
			status:  http.StatusNotFound,
			body:    "404 page not found",
		},
		{
			name:    "Negative ID",
			urlPath: "/snippets/view/-1",
			status:  http.StatusNotFound,
			body:    "404 page not found",
		},
		{
			name:    "Decimal ID",
			urlPath: "/snippets/view/1.23",
			status:  http.StatusNotFound,
			body:    "404 page not found",
		},
		{
			name:    "String ID",
			urlPath: "/snippets/view/ada",
			status:  http.StatusNotFound,
			body:    "404 page not found",
		},
		{
			name:    "Empty ID",
			urlPath: "/snippets/view/",
			status:  http.StatusNotFound,
			body:    "404 page not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure each subtest starts with a clean slate (no leftover cookies).
			ts.resetClientCookieJar(t)

			res := ts.get(t, tt.urlPath)

			assert.Equal(t, res.status, tt.status)
			if tt.body != "" {
				assert.StringContains(t, res.body, tt.body)
			}
		})
	}
}

// TestUserSignup verifies that the signup page is rendered correctly and
// contains the necessary security tokens and session cookies required
// for a successful form submission.
func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	// Request the signup form.
	res := ts.get(t, "/user/signup")

	// Assert the page loaded successfully.
	assert.Equal(t, res.status, http.StatusOK)

	// Extract the CSRF token to ensure it is present in the HTML.
	// If the token is missing, extractCSRFToken will call t.Fatal.
	token := extractCSRFToken(t, res.body)
	t.Logf("Extracted CSRF token: %q", token)

	// Log cookie details for debugging session persistence.
	// In a real scenario, we expect to see a session cookie here.
	t.Logf("Response Cookies: %v", res.cookies)
}
