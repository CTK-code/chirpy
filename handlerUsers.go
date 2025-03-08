package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CTK-code/chirpy/internal/auth"
	"github.com/CTK-code/chirpy/internal/database"
	"github.com/google/uuid"
)

type userRes struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding parameters", err)
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "Error hashing", err)
	}
	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, 400, "Error creating user", err)
	}

	response := userRes{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJson(w, 201, response)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, "incorrect email or password", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "incorrect email or password", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, 1*time.Hour)
	if err != nil {
		respondWithError(w, 400, "error creating jwt token", err)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 400, "error creating refresh token", err)
	}
	refreshParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}
	refresh, err := cfg.db.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, 400, "error creating refresh token", err)
	}
	response := userRes{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh.Token,
		IsChirpyRed:  user.IsChirpyRed,
	}
	respondWithJson(w, 200, response)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "error getting token", err)
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 401, "error unmarshaling parameters", err)
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "could not validate token", err)
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 401, "error hashing", err)
	}

	updateParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}
	user, err := cfg.db.UpdateUser(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, 401, "error updating user", err)
	}
	response := userRes{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJson(w, 200, response)
}
