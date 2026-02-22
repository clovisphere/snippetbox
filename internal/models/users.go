package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created_at)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	if _, err := m.DB.Exec(stmt, name, email, string(hashedPassword)); err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Authenticate verifies whether a user exists with the provided email address
// and password. If the user exists and the password matches the hash, it
// returns the user's ID.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	// Retrieve the id and hashed password associated with the given email.
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

// Exists checks if a user with a specific ID exists in the users table.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

// Get retrieves a specific user record based on its ID.
func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}
