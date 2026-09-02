package main

import (
	"net/http"
	"log"
)

func main() {
	// Create a ServeMux for routing requests
	serveMux := http.NewServeMux()
	
	// Add path handlers
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	serveMux.Handle("/assets", http.FileServer(http.Dir("./assets")))
	
	// Create a new http Server
	s := &http.Server{
		Addr:		":8080",
		Handler:	serveMux,
	}
	
	// Start the server
	log.Fatal(s.ListenAndServe())
}