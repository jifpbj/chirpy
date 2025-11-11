package main

import (
	"encoding/json"
	"net/http"

	"github.com/jifpbj/chirpy/internal/auth"
	"github.com/jifpbj/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't hash password", err)
		return
	}
	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: password,
	}

	user, err := apiCfg.db.CreateUser(req.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	out := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, out)
}
