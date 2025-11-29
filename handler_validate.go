package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type parameters struct {
	Body   string `json:"body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func cleanProfanity(body string) string {
	banned := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
	tokens := strings.Split(body, " ")
	for i, tok := range tokens {
		if banned[strings.ToLower(tok)] {
			tokens[i] = "****"
		}
	}
	return strings.Join(tokens, " ")
}
