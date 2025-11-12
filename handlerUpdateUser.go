package main

import (
	"encoding/json"
	"net/http"

	"github.com/jifpbj/chirpy/internal/auth"
	"github.com/jifpbj/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	BearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid token", err)
	}

	userID, err := auth.ValidateJWT(BearerToken, apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
	}

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't hash password", err)
		return

	}

	userParams := database.UpdateUserByIDParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: password,
	}

	user, err := apiCfg.db.UpdateUserByID(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
	}

	output := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusOK, output)
}
