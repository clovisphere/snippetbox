package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clovisphere/snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	// Setup the Recorder and Request
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a "next" handler that we can verify was called
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Run the middleware
	commonHeaders(next).ServeHTTP(rr, r)

	// Validate the Results
	res := rr.Result()
	defer res.Body.Close()

	// Use a table to verify all security headers in one go
	tests := []struct {
		headerName    string
		expectedValue string
	}{
		{"Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"},
		{"Referrer-Policy", "origin-when-cross-origin"},
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "deny"},
		{"X-XSS-Protection", "0"},
		{"Server", "Go"},
	}

	for _, tt := range tests {
		t.Run(tt.headerName, func(t *testing.T) {
			got := res.Header.Get(tt.headerName)
			if got != tt.expectedValue {
				t.Errorf("expected %s: %q; got %q", tt.headerName, tt.expectedValue, got)
			}
		})
	}

	// Verify the middleware actually passed execution to the next handler
	if !nextCalled {
		t.Error("next handler was not called")
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200; got %d", res.StatusCode)
	}

	// Read the body from the recorder
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the body was passed through correctly
	assert.Equal(t, string(bytes.TrimSpace(body)), "OK")
}
