package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/CTK-code/chirpy/internal/auth"
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
		Body string `json:"body"`
	}

	type res struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
		Token     string    `json:"token"`
	}

	// Check that the jwt token is valid
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error authorizing", err)
	}
	posterId, err := auth.ValidateJWT(token, conf.secret)
	if err != nil {
		respondWithError(w, 401, "error authorizing", err)
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding request", err)
		return
	}
	chirpParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: posterId,
	}
	ch, err := conf.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 400, "Error inserting chirp", err)
		return
	}
	response := res{
		Id:        ch.ID,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
		Body:      ch.Body,
		UserId:    ch.UserID,
		Token:     token,
	}
	respondWithJson(w, 201, response)
}

func (conf *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorIdQuery := r.URL.Query().Get("author_id")
	sortByQuery := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	if authorIdQuery != "" {
		authorId, err := uuid.Parse(authorIdQuery)
		if err != nil {
			respondWithError(w, 400, "Error getting chirps", err)
		}
		chirps, err = conf.db.GetChirpsByAuthor(r.Context(), authorId)
		if err != nil {
			respondWithError(w, 400, "could not find auhtors chirps", err)
		}
	} else {
		var err error
		chirps, err = conf.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, 400, "Error getting chirps", err)
		}

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
	if sortByQuery == "desc" {
		sort.Slice(chirpArr, func(i, j int) bool {
			return chirpArr[i].CreatedAt.After(chirpArr[j].CreatedAt)
		})
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

func (conf *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error parsing chirpID", err)
	}
	userId, err := auth.ValidateJWT(token, conf.secret)
	if err != nil {
		respondWithError(w, 401, "could not validate token", err)
	}
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Error parsing chirpID", err)
	}
	ch, err := conf.db.GetChirpById(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Could not retrieve chirp", err)
	}
	if ch.UserID != userId {
		w.WriteHeader(403)
		return
	}
	_, err = conf.db.DeleteChirp(r.Context(), ch.ID)
	if err != nil {
		respondWithError(w, 404, "Could not delete chirp", err)
	}
	w.WriteHeader(204)
}
