package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sync/atomic"
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/Ikit24/Chirpy/internal/database"

)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	const root = "."
	const port = "8080"

	apiCfg := apiConfig{db: dbQueries}
	mux := http.NewServeMux()

	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(root)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServer))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)

	srv := &http.Server{Addr: ":" + port, Handler: mux}
	log.Fatal(srv.ListenAndServe())
}
