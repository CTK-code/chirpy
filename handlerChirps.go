package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CTK-code/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

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
	ch, err := conf.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 400, "Error inserting chirp", err)
		return
	}
	response := chirp{
		Id:        ch.ID,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
		Body:      ch.Body,
		UserId:    ch.UserID,
	}
	respondWithJson(w, 201, response)
}

func (conf *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := conf.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "Error getting chirps", err)
	}
	chirpArr := []chirp{}
	for _, ch := range chirps {
		response := chirp{
			Id:        ch.ID,
			CreatedAt: ch.CreatedAt,
			UpdatedAt: ch.UpdatedAt,
			Body:      ch.Body,
			UserId:    ch.UserID,
		}
		chirpArr = append(chirpArr, response)
	}
	respondWithJson(w, 200, chirpArr)
}

func (conf *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Error parsing chirpID", err)
	}
	ch, err := conf.db.GetChirpById(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Could not retieve chirp", err)
	}
	response := chirp{
		Id:        ch.ID,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
		Body:      ch.Body,
		UserId:    ch.UserID,
	}
	respondWithJson(w, 200, response)
}
