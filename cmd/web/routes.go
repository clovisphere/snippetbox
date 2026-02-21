package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.index)
	mux.HandleFunc("GET /snippets", app.list)
	mux.HandleFunc("GET /snippets/{id}", app.show)
	mux.HandleFunc("POST /snippets", app.create)

	return mux
}
