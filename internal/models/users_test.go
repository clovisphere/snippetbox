package models

import (
	"testing"

	"github.com/clovisphere/snippetbox/internal/assert"
)

// TestUserModelExists validates the Exists method against the test database.
// It covers successful lookups, missing records, and edge cases.
func TestUserModelExists(t *testing.T) {
	// Skip the test if we are running in short mode (e.g., go test -short).
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	// Initialize the test database and the UserModel.
	db := newTestDB(t)
	m := UserModel{DB: db}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1, // Alice (from setup.sql)
			want:   true,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the method under test.
			exists, err := m.Exists(tt.userID)

			// Use your custom assertion helpers.
			assert.NilError(t, err)
			assert.Equal(t, exists, tt.want)
		})
	}
}
