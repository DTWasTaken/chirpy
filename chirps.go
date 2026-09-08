package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DTWasTaken/chirpy/internal/database"
	"github.com/google/uuid"
)

type createChirpRequest struct {
	Body		string `json:"body"`
	UserID		uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID			uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Body		string `json:"body"`
	UserID		uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestedChirp := createChirpRequest{}
	err := decoder.Decode(&requestedChirp)
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
	
	if len(requestedChirp.Body) > 140 {
		writeResponse(
			w,
			http.StatusBadRequest,
			errRespBody{
				"Chirp is too long",
			},
		)
		return
	}
	
	requestedChirp.Body = replaceBadWords(requestedChirp.Body)
	
	params := database.PostChirpParams{
		Body:	requestedChirp.Body,
		UserID:	requestedChirp.UserID,
	}
	
	newChirp, err := cfg.db.PostChirp(r.Context(), params)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			errRespBody{
				fmt.Sprintf("Error saving chirp: %s", err),
			},
		)
		return
	}

	writeResponse(
		w,
		http.StatusCreated,
		Chirp{
			ID:			newChirp.ID,
			CreatedAt:	newChirp.CreatedAt,
			UpdatedAt:	newChirp.UpdatedAt,
			Body:		newChirp.Body,
			UserID:		newChirp.UserID,
		},
	)
}

func replaceBadWords(input string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	
	var output []string
	
	for _, word := range strings.Split(input, " ") {
		if stringIsInStrings(strings.ToLower(word), badWords) {
			output = append(output, "****")
			continue
		}
		output = append(output, word)
	}
	
	return strings.Join(output, " ")
}

func stringIsInStrings(check string, strs []string) bool {
	for _, str := range strs {
		if str == check {
			return true
		}
	}
	return false
}