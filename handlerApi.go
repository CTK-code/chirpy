package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/CTK-code/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "malformated header", err)
	}
	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), authToken)
	if err != nil {
		respondWithError(w, 401, "refresh token not found", err)
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "refresh token has been revoked", errors.New("invalid token"))
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, 1*time.Hour)
	if err != nil {
		respondWithError(w, 400, "error creating jwt token", err)
	}

	respondWithJson(w, 200, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "malformated header", err)
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), authToken)
	if err != nil {
		respondWithError(w, 401, "token not found", err)
	}
	w.WriteHeader(204)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type valid struct {
		CreatedAt   time.Time `json:"created_at"`
		Valid       bool      `json:"valid"`
		CleanedBody string    `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	res := valid{
		CreatedAt:   time.Now(),
		Valid:       true,
		CleanedBody: filterBody(params.Body),
	}
	respondWithJson(w, http.StatusOK, res)
}

func filterBody(body string) string {
	cleanedBody := body
	words := strings.Split(cleanedBody, " ")
	for i, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[i] = censor
		}
	}
	return strings.Join(words, " ")
}

const censor = "****"

var profaneWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}
