package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type chirpPost struct {
	Body 	string `json:"body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	post := chirpPost{}
	err := decoder.Decode(&post)
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
	
	if len(post.Body) > 140 {
		writeResponse(
			w,
			http.StatusBadRequest,
			errRespBody{
				"Chirp is too long",
			},
		)
		return
	}
	
	post.Body = replaceBadWords(post.Body)
	
	writeResponse(
		w,
		http.StatusOK,
		cleanedRespBody{
			post.Body,
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