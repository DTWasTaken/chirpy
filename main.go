package main

import (
	"net/http"
	"log"
)

func main() {
	// Create a ServeMux for routing requests
	mux := http.NewServeMux()
	
	// Add path handlers
	mux.Handle(
		"/app/",
		http.StripPrefix(
			"/app/",
			http.FileServer(http.Dir(".")),
		),
	)
	mux.HandleFunc("/healthz", handlerHealthz)
	
	// Create a new http Server
	s := &http.Server{
		Addr:		":8080",
		Handler:	mux,
	}
	
	// Start the server
	log.Fatal(s.ListenAndServe())
}

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}