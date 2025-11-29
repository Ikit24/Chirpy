package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
	"sort"

	"github.com/google/uuid"
)

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	getChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve chirps")
		return
	}

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "bad id")
				return
			}
	}

	sortOrder := r.URL.Query().Get("sort")

	resp := make([]ChirpResponse, 0, len(getChirps))
	for _, c := range getChirps {
		if (authorID != uuid.Nil) && (c.UserID != authorID) {
			continue
		}
		resp = append(resp, ChirpResponse {
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	if sortOrder == "asc" || sortOrder == "" {
		sort.Slice(resp, func(i, j int) bool {
			return resp[i].CreatedAt.Before(resp[j].CreatedAt)
		})
	} else if sortOrder == "desc" {
		sort.Slice(resp, func(i, j int) bool {
			return resp[j].CreatedAt.Before(resp[i].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
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
	resp := ChirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
