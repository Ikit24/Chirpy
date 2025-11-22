package main

import (
	"database/sql"
	"net/http"
	"errors"

	"github.com/google/uuid"
	"github.com/Ikit24/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "bad id")
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate token")
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}
	
	if err := cfg.db.DeleteChirp(r.Context(), id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal error ")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
