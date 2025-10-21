package main

import (
	"encoding/json"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

type errResp struct {
	Error string `json:"error"`
}

type okResp struct {
	Valid bool `json:"valid"`
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		writeJSONError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	writeJSON(w, http.StatusOK, okResp{Valid: true})
}

func respondWithError(w http.RespondseWriter, code int, msg string) {}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {}
