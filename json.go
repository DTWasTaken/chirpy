package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type errRespBody struct {
	ErrorMessage	string `json:"error"`
}


func writeResponse(w http.ResponseWriter, code int, body any) {
	dat, err := json.Marshal(body)
	if err != nil {
		writeResponse(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error marshalling JSON: %s", err),
		)
		return
	}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}