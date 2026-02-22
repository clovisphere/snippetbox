package models

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql" // Blank import to register the MySQL driver
)

// newTestDB creates a connection to the test database, executes the setup.sql
// script to prepare the schema/data, and registers a cleanup function to
// run teardown.sql and close the connection when the test finishes.
func newTestDB(t *testing.T) *sql.DB {
	// Establish a connection to the test database on port 3307.
	// The multiStatements=true parameter is crucial here as it allows
	// our SQL driver to execute multiple commands (CREATE, INSERT) in one go.
	db, err := sql.Open("mysql", "test:test@tcp(localhost:3307)/snippetbox_test?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read and execute the setup script.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	if _, err := db.Exec(string(script)); err != nil {
		db.Close()
		t.Fatal(err)
	}

	// t.Cleanup registers a function to run after the current test (and all
	// its subtests) have completed. This ensures our database is always
	// wiped clean, even if a test panics.
	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := db.Exec(string(script)); err != nil {
			t.Fatal(err)
		}
	})

	return db
}
