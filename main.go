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
	secret         string
	polkaKey       string
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
	apiConf.secret = os.Getenv("SECRET")
	apiConf.polkaKey = os.Getenv("POLKA_KEY")

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConf.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiConf.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConf.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	mux.HandleFunc("POST /api/users", apiConf.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiConf.handlerLoginUser)
	mux.HandleFunc("POST /api/chirps", apiConf.handlerCreateChirp)
	mux.HandleFunc("POST /api/refresh", apiConf.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apiConf.revokeHandler)
	mux.HandleFunc("PUT /api/users", apiConf.handlerUpdateUser)

	mux.HandleFunc("GET /api/chirps", apiConf.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConf.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConf.handlerDeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", apiConf.handlerPolkaHook)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
