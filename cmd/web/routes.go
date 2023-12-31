package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("C:\\Users\\MSI\\Documents\\projects\\go\\lets-go-snippetbox\\ui\\static\\"))

	mux.Handle("C:\\Users\\MSI\\Documents\\projects\\go\\lets-go-snippetbox\\ui\\static\\", http.StripPrefix("\\static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
