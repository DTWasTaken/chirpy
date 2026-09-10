package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/DTWasTaken/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	db				*database.Queries
	platform		string
	clientSecret	string
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
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Database error: %s", err),
			},
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}