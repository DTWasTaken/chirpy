package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DTWasTaken/chirpy/internal/auth"
	"github.com/DTWasTaken/chirpy/internal/database"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
}

type User struct {
	ID				uuid.UUID `json:"id"`
	CreatedAt		time.Time `json:"created_at"`
	UpdatedAt		time.Time `json:"updated_at"`
	Email			string `json:"email"`
	AccessToken		string `json:"token"`
	RefreshToken	string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestedUser := createUserRequest{}
	err := decoder.Decode(&requestedUser)
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
	
	if requestedUser.Email == "" || requestedUser.Password == "" {
		writeResponse(
			w,
			http.StatusBadRequest,
			errRespBody{
				fmt.Sprintf("'email' and 'password' fields are required"),
			},
		)
		return
	}
	
	hashedPassword, err := auth.HashPassword(requestedUser.Password)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{err.Error()},
		)
		return
	}
	
	params := database.CreateUserParams{
		Email:			requestedUser.Email,
		HashedPassword:	hashedPassword,
	}

	newUser, err := cfg.db.CreateUser(r.Context(), params)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error creating user with email %s: %s",
							requestedUser.Email,
							err),
			},
		)
		return
	}

	writeResponse(
		w,
		http.StatusCreated,
		User{
			ID:			newUser.ID,
			CreatedAt:	newUser.CreatedAt,
			UpdatedAt:	newUser.UpdatedAt,
			Email:		newUser.Email,
		},
	)
}