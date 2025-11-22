package main

import(
	"net/http"
	"time"

	"github.com/Ikit24/Chirpy/internal/auth"
)

type respToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	newAccess, err := auth.MakeJWT(userID, cfg.secret, time.Hour)
		if err != nil{
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}
	out := respToken{Token: newAccess}
	respondWithJSON(w, http.StatusOK, out)
}
