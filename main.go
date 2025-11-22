package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Ikit24/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret		   string
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
    IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	sec := os.Getenv("SECRET")
	if sec == "" {
		log.Fatal("SECRET missing")
	}
	cfg := apiConfig{
		db:       queries,
		platform: os.Getenv("PLATFORM"),
		secret:   os.Getenv("SECRET"),
	}

	const root = "."
	const port = "8080"

	mux := http.NewServeMux()
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(root)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileServer))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/chirps", cfg.handlerReturnChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGet)

	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirps)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerWebhooks)

	mux.HandleFunc("PUT /api/users", cfg.handlerUsersPut)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerChirpsDelete)

	srv := &http.Server{Addr: ":" + port, Handler: mux}
	log.Fatal(srv.ListenAndServe())
}
