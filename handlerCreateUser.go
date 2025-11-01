package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"update_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		uuid.New(),
		time.Now(),
		time.Now(),
		params.Email,
	})
}
