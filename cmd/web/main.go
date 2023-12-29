package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("C:\\Users\\MSI\\Documents\\projects\\go\\lets-go-snippetbox\\ui\\static\\"))

	mux.Handle("C:\\Users\\MSI\\Documents\\projects\\go\\lets-go-snippetbox\\ui\\static\\", http.StripPrefix("\\static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
