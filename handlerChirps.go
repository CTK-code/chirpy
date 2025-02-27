package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CTK-code/chirpy/internal/database"
	"github.com/google/uuid"
)

func (conf *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type res struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding request", err)
		return
	}
	chirpParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserId,
	}
	chirp, err := conf.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 400, "Error inserting chirp", err)
		return
	}
	response := res{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}
	respondWithJson(w, 201, response)
}
