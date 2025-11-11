package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jifpbj/chirpy/internal/auth"
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
	output := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, output)
}
