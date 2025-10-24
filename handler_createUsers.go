package main

import (
	"encoding/json"
	"net/http"
	"log"
)

type userParams struct {
	Email string `json:"email"`
}

func respondError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	var params userParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Println("create error:", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	out := User{
	    ID:        dbUser.ID,
	    CreatedAt: dbUser.CreatedAt,
	    UpdatedAt: dbUser.UpdatedAt,
	    Email:     dbUser.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(out)
}

