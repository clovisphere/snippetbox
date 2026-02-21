package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *slog.Logger
}

// main parses flags, connects to the database, and starts the HTTP server.
func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "dev:demo@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: false,
				Level:     slog.LevelDebug,
			},
		),
	)

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error("Could not connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	// Initialize the application
	app := &application{
		logger: logger,
	}

	logger.Info("Starting server", slog.String("addr", *addr))

	if err := http.ListenAndServe(*addr, app.routes()); !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Server gracefully shutdown 😊")
}

// openDB opens a connection to the MySQL database using the given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}
