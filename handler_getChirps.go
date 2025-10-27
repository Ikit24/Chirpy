package main

import (
	"net/http"
	"log"
)

func (cfg *apiConfig) handlerReturnChirps(w http.ResponseWriter, r *http.Request) {
	getChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error retrieving chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve chirps")
		return
	}
	resp := make([]ChirpResponse, 0, len(getChirps))
	for _, c := range getChirps {
		resp = append(resp, ChirpResponse{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}
