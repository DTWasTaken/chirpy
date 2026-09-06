package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/google/uuid"
)

type createUserRequest struct {
	Email 	string `json:"email"`
}

type User struct {
	ID			uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Email		string `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
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

	newUser, err := cfg.db.CreateUser(r.Context(), requestedUser.Email)
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