package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ikit24/Chirpy/internal/auth"
	"github.com/Ikit24/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginParams struct {
        Email    string `json:"email"`
        Password string `json:"password"`
	}

	var params loginParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	dbLogin, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, dbLogin.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token")
		return
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour)
	err = cfg.db.InsertRefreshToken(r.Context(), database.InsertRefreshTokenParams {
		Token: refreshToken,
		UserID: dbLogin.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token")
		return
	}

	token, err := auth.MakeJWT(dbLogin.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create token")
		return
	}

	out := User{
		ID:        dbLogin.ID,
		CreatedAt: dbLogin.CreatedAt,
		UpdatedAt: dbLogin.UpdatedAt,
		Email:     dbLogin.Email,
		Token:     token,
		RefreshToken: refreshToken,
		IsChirpyRed: dbLogin.IsChirpyRed,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}
