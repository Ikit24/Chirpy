package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ikit24/Chirpy/internal/auth"
	"github.com/Ikit24/Chirpy/internal/database"
)

type updateUserParams struct {
    Email    	string `json:"email"`
    Password 	string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerPutUsers(w http.ResponseWriter, r *http.Request) {
	var params updateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params) ; err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	dbUser, err := cfg.db.UpdateUserById(
		r.Context(),
		database.UpdateUserByIdParams{
			ID: 			userID,
			Email: 			params.Email,
			HashedPassword: hash,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update user")
		return
	}

	out := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}
