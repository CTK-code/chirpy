package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"
)

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
