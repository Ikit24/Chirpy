package main

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/Ikit24/Chirpy/internal/auth"
	"github.com/Ikit24/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type createUserParams struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
	}

	var params createUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	out := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(out)
}
