package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DTWasTaken/chirpy/internal/auth"
	"github.com/DTWasTaken/chirpy/internal/database"
	"github.com/google/uuid"
)

type polkaWebhookRequest struct {
	Event		string `json:"event"`
	Data		struct{
		UserID		uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			writeResponse(
				w,
				http.StatusUnauthorized,
				errRespBody{"Invalid Authorization header"},
			)
			return
		}
	
	if apiKey != cfg.polkaAPIKey {
		writeResponse(
			w,
			http.StatusUnauthorized,
			errRespBody{"Invalid API key"},
		)
		return
	}
	
	decoder := json.NewDecoder(r.Body)
	webhook := polkaWebhookRequest{}
	err = decoder.Decode(&webhook)
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
	
	if webhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	
	params := database.SetUserRedStatusParams{
		ID:				webhook.Data.UserID,
		IsChirpyRed:	true,
	}
	
	_, err = cfg.db.SetUserRedStatus(r.Context(), params)
	if err != nil {
		writeResponse(
			w,
			http.StatusNotFound,
			errRespBody{"Could not find user"},
		)
		return
	}
		
	w.WriteHeader(http.StatusNoContent)
}