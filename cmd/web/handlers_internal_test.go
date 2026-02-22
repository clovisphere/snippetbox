package main

import (
	"net/http"
	"strings"
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
	const expectedString = "Latest Snippets"
	if !strings.Contains(res.body, expectedString) {
		t.Errorf("expected body to contain %q", expectedString)
	}
}
