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
	dbQueries      *database.Queries
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
	apiConf.dbQueries = database.New(db)
	mux := http.NewServeMux()
	mux.Handle("/app/", apiConf.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiConf.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConf.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (conf *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	conf.fileserverHits.Store(0)
}
