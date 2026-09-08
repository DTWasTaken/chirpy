package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DTWasTaken/chirpy/internal/auth"
)

type loginRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
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
	
	writeResponse(
		w,
		http.StatusOK,
		User{
			ID:			user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt:	user.UpdatedAt,
			Email:		user.Email,
		},
	)
}