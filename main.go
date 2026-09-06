package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/DTWasTaken/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file for postgres connection string
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	
	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open database")
	}
	
	dbQueries := database.New(db)
	
	// Create a ServeMux for routing requests
	mux := http.NewServeMux()
	
	apiCfg := apiConfig{
		fileserverHits:	atomic.Int32{},
		dbQueries:		dbQueries,
	}
	
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

