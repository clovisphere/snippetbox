package models

import (
	"database/sql"
	"errors"
	"time"
)

// Storage provides methods to interact with the database tables.
type Storage struct {
	DB *sql.DB
}

// Snippet represents a single snippet record in the database.
type Snippet struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Insert adds a new snippet to the database.
// title: snippet title, content: snippet body, expires: days until expiration.
// Returns the ID of the new snippet or an error.
func (s *Storage) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created_at, expires_at)
	VALUES(?, ?, UTC_TIMESTAMP, DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := s.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get retrieves a snippet by ID.
// Returns an error if no matching snippet is found.
func (s *Storage) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created_at, expires_at FROM snippets
	WHERE expires_at > UTC_TIMESTAMP() AND id = ?`

	row := s.DB.QueryRow(stmt, id)

	var snippet Snippet
	if err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.CreatedAt, &snippet.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return snippet, nil
}

// Latest returns the 10 most recently created, non-expired snippets.
func (s *Storage) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created_at, expires_at FROM snippets
	WHERE expires_at > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet
	for rows.Next() {
		var snippet Snippet
		if err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.CreatedAt, &snippet.ExpiresAt); err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
