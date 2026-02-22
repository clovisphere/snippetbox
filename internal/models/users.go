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

// UserModelInterface defines the set of methods required for user-related
// data operations. This abstraction enables the application to switch
// between a production database (UserModel) and a mock implementation (mocks.UserModel)
// for testing authentication and authorization flows.
type UserModelInterface interface {
	// Insert adds a new user record. It should return ErrDuplicateEmail if
	// the email address is already in use.
	Insert(name, email, password string) error

	// Authenticate verifies if a user exists with the provided credentials.
	// It returns the user's ID on success, or ErrInvalidCredentials on failure.
	Authenticate(email, password string) (int, error)

	// Exists checks if a specific user ID exists in the system.
	Exists(id int) (bool, error)

	// Get retrieves a full User record by its unique ID.
	Get(id int) (*User, error)
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

// Exists checks if a user with a specific ID exists in the database.
// It returns true if the user exists, otherwise false.
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	// Use the EXISTS() function which returns 1 (true) or 0 (false).
	// This is more performant than selecting a count or the full record.
	query := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(query, id).Scan(&exists)

	return exists, err
}

// Get retrieves a specific user record based on its ID.
func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}
