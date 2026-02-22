package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/clovisphere/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// application holds the dependencies for the web application, including
// a form decder, a logger, a session manager, the storage layer for database access, and the template cache.
type application struct {
	formDecoder    *form.Decoder
	logger         *slog.Logger
	sessionManager *scs.SessionManager
	storage        *models.Storage
	templateCache  map[string]*template.Template
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

	// Connect to the MySQL database using the provided DSN
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error("Could not connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	// Build a template cache by parsing all page templates along with
	// the base layout and partials
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error("Could not load template cache", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize a new form decoder instance to map form data onto Go structs
	formDecoder := form.NewDecoder()

	// Initialize a new session manager using a MySQL-backed storage engine
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize the application struct with form decoder, logger, session manager, storage layer, and templates
	app := &application{
		formDecoder:    formDecoder,
		logger:         logger,
		sessionManager: sessionManager,
		storage:        &models.Storage{DB: db},
		templateCache:  templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting server", slog.String("addr", srv.Addr))

	// Start the HTTP server and listen on the specified address
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
