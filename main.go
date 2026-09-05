package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a ServeMux for routing requests
	mux := http.NewServeMux()
	
	apiCfg := apiConfig{}
	
	// Add path handlers
	// /app
	mux.Handle(
		"/app/",
		http.StripPrefix(
			"/app/",
			apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("."))),
		),
	)
	
	// /api
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	
	// /admin
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	
	// Create a new http Server
	s := &http.Server{
		Addr:		":8080",
		Handler:	mux,
	}
	
	// Start the server
	log.Fatal(s.ListenAndServe())
}

