package main

import (
	"net/http"
	"net/url"
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

// TestUserSignup performs a suite of table-driven tests for the user registration
// flow. It validates successful creation, validation errors, and CSRF protection.
func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	// Define standard test data constants.
	const (
		validName     = "Ken Thompson"
		validPassword = "@V3ry$3cur3P@$$w0rd!"
		validEmail    = "ken@example.com"
		signupFormTag = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct {
		description       string
		name              string
		email             string
		password          string
		useValidCSRFToken bool
		expectedStatus    int
		expectedBody      string
	}{
		{
			description:       "Valid submission",
			name:              validName,
			email:             validEmail,
			password:          validPassword,
			useValidCSRFToken: true,
			expectedStatus:    http.StatusSeeOther, // Redirects to login on success
		},
		{
			description:       "Missing CSRF token",
			name:              validName,
			email:             validEmail,
			password:          validPassword,
			useValidCSRFToken: false,
			expectedStatus:    http.StatusBadRequest,
		},
		{
			description:       "Invalid email format",
			name:              validName,
			email:             "not-an-email",
			password:          validPassword,
			useValidCSRFToken: true,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedBody:      signupFormTag,
		},
		{
			description:       "Short password",
			name:              validName,
			email:             validEmail,
			password:          "123",
			useValidCSRFToken: true,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedBody:      signupFormTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			// GET the form to obtain the CSRF token
			res := ts.get(t, "/user/signup")

			form := url.Values{}
			form.Add("name", tt.name)
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			if tt.useValidCSRFToken {
				form.Add("csrf_token", extractCSRFToken(t, res.body))
			}

			// POST the data
			res = ts.postForm(t, "/user/signup", form)

			// Assert status and body content
			assert.Equal(t, res.status, tt.expectedStatus)
			if tt.expectedBody != "" {
				assert.StringContains(t, res.body, tt.expectedBody)
			}
		})
	}
}
