package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/Ikit24/Chirpy/internal/auth"
)

type polkaWebhookData struct {
	UserID string `json:"user_id"`
}

type polkaWebhookRequest struct {
	Event string		   `json:"event"`
	Data  polkaWebhookData `json:"data"`
}

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid API key")
		return
	}

	var polkaData polkaWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&polkaData) ; err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if polkaData.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id, err := uuid.Parse(polkaData.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
