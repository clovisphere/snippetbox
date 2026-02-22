// Package assert provides basic assertion helpers for unit tests to reduce
// boilerplate and improve test readability.
package assert

import "testing"

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
