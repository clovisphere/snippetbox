package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testServer wraps a httptest.Server instance to provide custom helper methods.
type testServer struct {
	*httptest.Server
}

// testResponse represents the essential parts of an HTTP response
// captured during testing for easier assertions.
type testResponse struct {
	status  int
	headers http.Header
	cookies []*http.Cookie
	body    string
}

// newTestApplication initializes a minimal application instance with a discarded
// logger to keep test output clean.
func newTestApplication(t *testing.T) *application {
	return &application{
		// slog.DiscardHandler prevents logs from cluttering the 'go test' output.
		logger: slog.New(slog.DiscardHandler),
	}
}

// newTestServer starts a TLS test server with the provided handler.
// It uses t.Cleanup to automatically shut down the server when the test finishes.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// Automatically close the server when the test (and sub-tests) complete.
	t.Cleanup(func() {
		ts.Close()
	})

	return &testServer{ts}
}

// get performs a GET request against the test server for a given path.
// It returns a testResponse containing the status, headers, cookies, and body.
func (ts *testServer) get(t *testing.T, urlPath string) testResponse {
	// The ts.Client() is pre-configured to trust the TLS certificate of the test server.
	res, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}
