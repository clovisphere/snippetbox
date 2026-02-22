package models

import (
	"database/sql"
	"time"
)

// User represents the data held in the users table for an individual user.
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
}

// UserModel wraps a sql.DB connection pool and provides methods
// for interacting with the users table.
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the users table. If the email already exists,
// it returns an ErrDuplicateEmail error.
func (s *UserModel) Insert(name, email, password string) error {
	// TODO: Use bcrypt to hash the password.
	// TODO: Execute a SQL INSERT statement.
	// TODO: Handle the MySQL 1062 error for duplicate emails.
	return nil
}

// Authenticate verifies whether a user exists with the provided email address
// and password. If the user exists and the password is correct, it returns
// the user's ID.
func (s *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists checks if a user with a specific ID exists in the users table.
func (s *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

// Get retrieves a specific user record based on its ID.
func (s *UserModel) Get(id int) (*User, error) {
	return nil, nil
}
