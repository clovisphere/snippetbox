package mocks

import (
	"time"

	"github.com/clovisphere/snippetbox/internal/models"
)

// mockUser represents a standard "active" user for testing purposes.
// The Name and Email honor the first computer programmer, Ada Lovelace.
// The HashedPassword corresponds to the string "pa$$word".
var mockUser = &models.User{
	ID:             1,
	Name:           "Ada Lovelace",
	Email:          "ada@example.com",
	HashedPassword: []byte("$2a$12$NuTj99w8Oc7otGu6TurtuO.p7q.E8T9v0C3sk9.94K3mQd.YqG/ly"),
	CreatedAt:      time.Now(),
}

// UserModel mocks the internal/models.UserModelInterface.
type UserModel struct{}

// Insert simulates creating a new user.
// It returns models.ErrDuplicateEmail if the email "dupe@example.com" is provided,
// simulating a database unique constraint violation.
func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

// Authenticate simulates the credential verification process.
// - Returns ID 1 if email and password match our mockUser ("alice@example.com" / "pa$$word").
// - Returns models.ErrInvalidCredentials for any other combination.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "ada@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

// Exists simulates checking for a user's presence in the database.
// It returns true only for ID 1.
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

// Get simulates retrieving a user's full profile.
// It returns the mockUser for ID 1, and models.ErrNoRecord for anything else.
func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
