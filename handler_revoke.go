package main

import(
	"net/http"

	"github.com/Ikit24/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	err = cfg.db.UpdateRevokedRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}
