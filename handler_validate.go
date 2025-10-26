package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type parameters struct {
	Body string `json:"body"`
	User_ID string `json:"user_id"`
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

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := cleanProfanity(params.Body)

	dbChirp, err := cfg.db.CreateChirp(r.Context(), params.Body, Parse(params.User_ID))
	if err != nil {
		log.Println("chirp error:", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	chrp := Chirp{
		ID:	   dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:	   dbChirp.Body,
		User_ID:   dbChirp.User_ID
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(chrp)
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
