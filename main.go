package main

import (
	"log"
	"net/http"
	"fmt"
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
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

func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("GET /api/healthz", apiHandler{})
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handlerReset)
	
	log.Fatal(server.ListenAndServe())
}
