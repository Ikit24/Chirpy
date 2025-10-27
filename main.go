package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sync/atomic"
	"database/sql"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/Ikit24/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	cfg := apiConfig{
		db:	queries,
		platform: os.Getenv("PLATFORM"),
	}
	
	const root = "."
	const port = "8080"

	mux := http.NewServeMux()
	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(root)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileServer))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/chirps", cfg.handlerReturnChirps)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirps)
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)

	srv := &http.Server{Addr: ":" + port, Handler: mux}
	log.Fatal(srv.ListenAndServe())
}
