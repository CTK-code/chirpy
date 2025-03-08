package main

import (
	"encoding/json"
	"net/http"

	"github.com/CTK-code/chirpy/internal/database"
	"github.com/google/uuid"
)

func (conf *apiConfig) handlerPolkaHook(w http.ResponseWriter, r *http.Request) {
	const expected = "user.upgraded"

	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var req request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 400, "error decoding json", err)
	}

	if req.Event != expected {
		w.WriteHeader(204)
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, 400, "error parsing id", err)
	}

	args := database.UpdateIsChirpyRedByIdParams{
		IsChirpyRed: true,
		ID:          userID,
	}
	_, err = conf.db.UpdateIsChirpyRedById(r.Context(), args)
	if err != nil {
		respondWithError(w, 404, "could not find user", err)
	}

	w.WriteHeader(204)
}
