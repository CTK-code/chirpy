package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/CTK-code/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	platform       string
	db             *database.Queries
}

func (conf *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incrementing metrics")
		conf.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiConf := apiConfig{}
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Connection to the DB failed")
		return
	}
	apiConf.db = database.New(db)
	apiConf.platform = os.Getenv("PLATFORM")
	mux := http.NewServeMux()
	mux.Handle("/app/", apiConf.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiConf.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConf.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	mux.HandleFunc("POST /api/users", apiConf.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiConf.handlerLoginUser)
	mux.HandleFunc("POST /api/chirps", apiConf.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConf.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConf.handlerGetChirp)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
