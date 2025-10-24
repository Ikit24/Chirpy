package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := cfg.db.Reset(r.Context()); err != nil {
		log.Println("reset error:", err)
		respondWithError(w, http.StatusInternalServerError, "reset failed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
