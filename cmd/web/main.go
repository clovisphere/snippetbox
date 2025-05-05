package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

// main entry point of the application.
func main() {
	os.Exit(start())
}

// start initializes the application and starts the server.
func start() int {
	addr := flag.String("addr", ":4000", "HTTP network address")
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

	// Initialize the application
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	return 1 // TODO: To be refactored
}
