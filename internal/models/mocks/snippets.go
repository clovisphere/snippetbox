// Package mocks provides mock implementations of the application's data models
// for use in unit and integration testing.
package mocks

import (
	"time"

	"github.com/clovisphere/snippetbox/internal/models"
)

// mockSnippet is a pre-defined Snippet used as a consistent return value
// for mock methods to ensure predictable test outcomes.
var mockSnippet = models.Snippet{
	ID:        1,
	Title:     "Gopher Standard Library Cheat Sheet",
	Content:   "Use the 'net/http' package for web servers and 'encoding/json' for API responses.",
	CreatedAt: time.Now(),
	ExpiresAt: time.Now(),
}

// SnippetModel mocks the internal/models.SnippetModelInterface.
// It provides hardcoded behaviors to simulate database interactions without
// requiring a live database connection.
type SnippetModel struct{}

// Insert simulates adding a new snippet to the database.
// It always returns a fixed ID of 2 and a nil error to simulate a successful creation.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	return 2, nil
}

// Get simulates fetching a single snippet by its ID.
// - If the ID is 1, it returns the mockSnippet.
// - For any other ID, it returns models.ErrNoRecord to simulate a "Not Found" state.
func (m *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		// Return an empty Snippet and the standard error for missing records
		return models.Snippet{}, models.ErrNoRecord
	}
}

// Latest simulates fetching the most recently created snippets.
// It returns a slice containing only the mockSnippet.
func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
