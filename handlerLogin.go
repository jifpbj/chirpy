package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jifpbj/chirpy/internal/auth"
	"github.com/jifpbj/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
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
	user, err := apiCfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user by email", err)
		return
	}
	if check, err := auth.CheckPasswordHash(params.Password, user.HashedPassword); !check {
		fmt.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	TokenExpirationDuration := time.Duration(3600) * time.Second

	token, err := auth.MakeJWT(user.ID, apiCfg.secret, TokenExpirationDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
	}

	apiCfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})

	output := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, output)
}
