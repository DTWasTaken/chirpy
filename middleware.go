package main

import (
	"net/http"

	"github.com/DTWasTaken/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareAuthenticated(
	handler func(w http.ResponseWriter, r *http.Request, userID uuid.UUID),
) func(http.ResponseWriter, *http.Request) {
	
	return func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			writeResponse(
				w,
				http.StatusUnauthorized,
				errRespBody{"Invalid Authorization header"},
			)
			return
		}
		
		userID, err := auth.ValidateJWT(bearerToken, cfg.clientSecret)
		if err != nil {
			writeResponse(
				w,
				http.StatusUnauthorized,
				errRespBody{"Invalid Authorization token"},
			)
			return
		}
		
		if userID == uuid.Nil {
			writeResponse(
				w,
				http.StatusUnauthorized,
				errRespBody{"Unauthorized"},
			)
			return
		}
		
		handler(w, r, userID)
	}
}