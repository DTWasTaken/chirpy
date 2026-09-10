package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/DTWasTaken/chirpy/internal/auth"
	"github.com/DTWasTaken/chirpy/internal/database"
	"github.com/google/uuid"
)

type userRequest struct {
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
	email, hashedPassword, errCode, err := getEmailAndHashedPasswordFromBody(r.Body)
	if err != nil {
		writeResponse(
			w,
			errCode,
			errRespBody{err.Error()},
		)
		return
	}
	
	params := database.CreateUserParams{
		Email:			email,
		HashedPassword:	hashedPassword,
	}

	newUser, err := cfg.db.CreateUser(r.Context(), params)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error creating user with email %s: %s", email, err),
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

func getEmailAndHashedPasswordFromBody(body io.ReadCloser) (email string, hashedPassword string, errCode int, err error) {
	decoder := json.NewDecoder(body)
	user := userRequest{}
	jsonErr := decoder.Decode(&user)
	if jsonErr != nil {
		return "", "", http.StatusInternalServerError, jsonErr
	}
	
	if user.Email == "" || user.Password == "" {
		return "", "", http.StatusBadRequest, fmt.Errorf("'email' and 'password' fields are required")
	}
	
	hashedPassword, hashErr := auth.HashPassword(user.Password)
	if hashErr != nil {
		return "", "", http.StatusInternalServerError, hashErr
	}
	
	return user.Email, hashedPassword, 0, nil
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	email, hashedPassword, errCode, err := getEmailAndHashedPasswordFromBody(r.Body)
	if err != nil {
		writeResponse(
			w,
			errCode,
			errRespBody{err.Error()},
		)
		return
	}
	
	params := database.UpdateUserParams{
		ID:				userID,
		Email:			email,
		HashedPassword:	hashedPassword,
	}
	
	updatedUser, err := cfg.db.UpdateUser(r.Context(), params)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error updating user with email %s: %s", email, err),
			},
		)
		return
	}

	writeResponse(
		w,
		http.StatusOK,
		User{
			ID:			updatedUser.ID,
			CreatedAt:	updatedUser.CreatedAt,
			UpdatedAt:	updatedUser.UpdatedAt,
			Email:		updatedUser.Email,
		},
	)
}