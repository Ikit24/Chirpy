package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiHandler struct{}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf(`
	<html>
	  <body>
	    <h1>Welcome, Chirpy Admin</h1>
	    <p>Chirpy has been visited %d times!</p>
	  </body>
	</html>`, hits)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

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
	w.Header().Set("Content-Type", "application/json")

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		resp, _ := json.Marshal(errResp{Error: "Something went wrong"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	if len(params.Body) > 140 {
		resp, _ := json.Marshal(errResp{Error: "Chirp is too long"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(okResp{Valid: true})
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("GET /api/healthz", apiHandler{})
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)

	log.Fatal(server.ListenAndServe())
}
