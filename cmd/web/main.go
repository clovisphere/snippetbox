package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
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

	logger.Info("Starting server", slog.String("addr", *addr))

	if err := http.ListenAndServe(*addr, app.routes()); !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
	}

	logger.Info("Server gracefully shutdown 😊")
}
