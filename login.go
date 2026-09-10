package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DTWasTaken/chirpy/internal/auth"
	"github.com/DTWasTaken/chirpy/internal/database"
)

type loginRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
}

type tokenResponse struct {
	Token	string `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	login := loginRequest{}
	err := decoder.Decode(&login)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error unmarshalling JSON: %s", err),
			},
		)
		return
	}
	
	user, err := cfg.db.GetUserByEmail(r.Context(), login.Email)
	if err != nil {
		writeResponse(
			w,
			http.StatusUnauthorized,
			errRespBody{"Incorrect email or password"},
		)
		return
	}
	
	match, err := auth.CheckPasswordHash(login.Password, user.HashedPassword)
	if !match || err != nil {
		writeResponse(
			w,
			http.StatusUnauthorized,
			errRespBody{"Incorrect email or password"},
		)
		return
	}
	
	expiresIn := time.Duration(time.Second * 3600)
	accessToken, err := auth.MakeJWT(user.ID, cfg.clientSecret, expiresIn)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error creating token: %s", err),
			},
		)
		return
	}
	
	refreshToken := auth.MakeRefreshToken()
	
	params := database.StoreRefreshTokenParams{
		Token:		refreshToken,
		ExpiresAt:	time.Now().UTC().Add(time.Hour * 24 * 60),
		UserID:		user.ID,
	}
	
	_, err = cfg.db.StoreRefreshToken(r.Context(), params)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error storing refresh token: %s", err),
			},
		)
		return
	}
	
	writeResponse(
		w,
		http.StatusOK,
		User{
			ID:				user.ID,
			CreatedAt:		user.CreatedAt,
			UpdatedAt:		user.UpdatedAt,
			Email:			user.Email,
			AccessToken:	accessToken,
			RefreshToken:	refreshToken,
		},
	)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeResponse(
			w,
			http.StatusUnauthorized,
			errRespBody{"Invalid Authorization header"},
		)
		return
	}
	
	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		writeResponse(
			w,
			http.StatusUnauthorized,
			errRespBody{"Invalid Token"},
		)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		writeResponse(
			w,
			http.StatusUnauthorized,
			errRespBody{"Token is expired"},
		)
		return
	}
	if refreshToken.RevokedAt.Valid {
		writeResponse(
			w,
			http.StatusUnauthorized,
			errRespBody{"Token is revoked"},
		)
		return
	}
	
	expiresIn := time.Duration(time.Second * 3600)
	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.clientSecret, expiresIn)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error creating token: %s", err),
			},
		)
		return
	}
	
	writeResponse(
		w,
		http.StatusOK,
		tokenResponse{accessToken},
	)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeResponse(
			w,
			http.StatusUnauthorized,
			errRespBody{"Invalid Authorization header"},
		)
		return
	}
	
	_, err = cfg.db.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error revoking token: %s", err),
			},
		)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}