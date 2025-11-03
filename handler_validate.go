package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Ikit24/Chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
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

	cleanedBody := cleanProfanity(params.Body)

	parsedUserID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID format")
		return
	}

	createChirpParams := database.CreateChirpsParams{
		Body:   cleanedBody,
		UserID: parsedUserID,
	}

	dbChirp, err := cfg.db.CreateChirps(r.Context(), createChirpParams)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}
	responseChirp := ChirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, responseChirp)
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
