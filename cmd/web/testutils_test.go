package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
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

// newTestServer initializes and returns a new TLS test server.
// The server is configured with a cookie jar to support session persistence
// and a custom redirect policy to allow inspection of intermediate responses.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize the TLS server with the provided handler.
	ts := httptest.NewTLSServer(h)

	// Initialize a cookie jar. This allows the test client to automatically
	// store and send cookies (like session IDs) across multiple requests.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server's client.
	ts.Client().Jar = jar

	// Disable automatic redirect following. By returning http.ErrUseLastResponse,
	// the client will stop and return the 302/303 response itself, allowing
	// us to assert on redirect locations and headers.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// Register a cleanup function to automatically shut down the server
	// when the test (and any sub-tests) finishes.
	t.Cleanup(func() {
		ts.Close()
	})

	return &testServer{ts}
}

// resetClientCookieJar replaces the existing cookie jar in the test server's client
// with a new, empty one. This effectively clears all cookies (including session
// and CSRF cookies) to simulate a fresh browser state or a logout.
func (ts *testServer) resetClientCookieJar(t *testing.T) {
	// Initialize a new, empty cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Update the test server's client with the new jar.
	// The old jar and its cookies will be garbage collected.
	ts.Client().Jar = jar
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
