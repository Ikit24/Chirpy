package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token")
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
