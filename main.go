package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type response struct {
	StatusCode	int
	Body		respBody
}

func (r response) Write(w http.ResponseWriter) {
	dat, err := r.Body.GetBodyJSON()
	if err != nil {
		resp := response{
			StatusCode: http.StatusInternalServerError,
			Body: errRespBody{
				ErrorMessage: fmt.Sprintf("Error marshalling JSON: %s", err),
			},
		}
		resp.Write(w)
		return
	}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(r.StatusCode)
    w.Write(dat)
}

type respBody interface {
	GetBodyJSON() ([]byte, error)
}

type errRespBody struct {
	ErrorMessage	string `json:"error"`
}

func (e errRespBody) GetBodyJSON() ([]byte, error) {
	return json.Marshal(e)
}

type validRespBody struct {
	Valid	bool `json:"valid"`
}

func (v validRespBody) GetBodyJSON() ([]byte, error) {
	return json.Marshal(v)
}

type chirpPost struct {
	Body 	string `json:"body"`
}

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

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type apiConfig struct {
	fileserverHits	atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	content := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, cfg.fileserverHits.Load())
	w.Write([]byte(content))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	post := chirpPost{}
	err := decoder.Decode(&post)
	if err != nil {
		resp := response{
			StatusCode: http.StatusInternalServerError,
			Body: errRespBody{
				ErrorMessage: fmt.Sprintf("Error unmarshalling JSON: %s", err),
			},
		}
		resp.Write(w)
		return
	}
	
	if len(post.Body) > 140 {
		resp := response{
			StatusCode: http.StatusBadRequest,
			Body: errRespBody{
				ErrorMessage: fmt.Sprintf("Chirp is too long"),
			},
		}
		resp.Write(w)
		return
	}
	
	resp := response{
		StatusCode: http.StatusOK,
		Body: validRespBody{
			Valid: true,
		},
	}
	resp.Write(w)
}