// Package assert provides basic assertion helpers for unit tests to reduce
// boilerplate and improve test readability.
package assert

import (
	"strings"
	"testing"
)

// Equal compares two comparable values of the same type.
// If the values are not equal, it fails the test and logs the mismatch.
// It uses the T.Helper() mechanism to ensure error reports point to
// the line where Equal was called, rather than inside this utility.
func Equal[T comparable](t *testing.T, actual, expected T) {
	// Marks this function as a test helper so the failure line
	// number points to the caller's location.
	t.Helper()

	if actual != expected {
		// Using quoted formatting (%v or %#v) helps identify empty strings
		// or specific types during a failure.
		t.Errorf("\n[Assert Equal Failed]\n  actual:   %v\n  expected: %v", actual, expected)
	}
}

// StringContains is a test helper that verifies if a string (actual) contains
// a specific substring (expectedSubstring). If the substring is not found,
// it fails the test with a descriptive error message.
func StringContains(t *testing.T, actual, expectedSubstring string) {
	// t.Helper() marks this function as a test helper, ensuring that
	// when a test fails, the line number reported points to the
	// calling test function rather than this helper.
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

// NilError is a test helper that checks if an error is nil.
// If it is not, it reports an error and identifies the calling
// line as the source of the failure.
func NilError(t *testing.T, actual error) {
	// t.Helper() ensures that the failure is reported at the
	// line where NilError was called, rather than inside this function.
	t.Helper()

	if actual != nil {
		t.Errorf("expected: nil error; got: %v", actual)
	}
}
