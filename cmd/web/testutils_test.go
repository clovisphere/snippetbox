package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/clovisphere/snippetbox/internal/models/mocks"
	"github.com/go-playground/form/v4"
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

// csrfTokenRX is a pre-compiled regular expression that captures the value of
// the CSRF token from an HTML hidden input field. We compile it once at
// package level for better performance during test execution.
var csrfTokenRX = regexp.MustCompile(`<input type=['"]hidden['"] name=['"]csrf_token['"] value=['"](.+)['"]\s*/?>`)

// newTestApplication initializes a minimal application instance with a discarded
// logger to keep test output clean.
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		// slog.DiscardHandler prevents logs from cluttering the 'go test' output.
		logger:         slog.New(slog.DiscardHandler),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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

// extractCSRFToken is a test helper that retrieves the CSRF token value from
// an HTML response body. It fails the test immediately if no token is found.
func extractCSRFToken(t *testing.T, body string) string {
	// t.Helper() marks this function as a test helper so that error reports
	// point to the actual test line rather than this function.
	t.Helper()

	// FindStringSubmatch returns a slice containing the full match and the
	// captured group. We expect a length of 2: [Full Match, Captured Token].
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	// Unescape the string to handle characters like '+' or '=' correctly
	// if they were encoded in the HTML.
	return html.UnescapeString(matches[1])
}
