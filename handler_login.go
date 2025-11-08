package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ikit24/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var params userParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	var expirationSec int

	if params.ExpiresInSeconds == 0 {
		expirationSec = 3600
	} else if params.ExpiresInSeconds > 3600 {
		expirationSec = 3600
	} else {
		expirationSec = params.ExpiresInSeconds
	}

	dbLogin, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, dbLogin.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(dbLogin.ID, cfg.secret, time.Duration(expirationSec) * time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}

	out := User{
		ID:        dbLogin.ID,
		CreatedAt: dbLogin.CreatedAt,
		UpdatedAt: dbLogin.UpdatedAt,
		Email:     dbLogin.Email,
		Token:     token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}
