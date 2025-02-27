package main

import (
	"fmt"
	"net/http"
)

func (conf *apiConfig) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,
		conf.fileserverHits.Load())
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (conf *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if conf.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("Unauthorized"))
		return
	}
	conf.fileserverHits.Store(0)
	conf.db.DeleteUsers(r.Context())
	w.WriteHeader(http.StatusOK)
}
