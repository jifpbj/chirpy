package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse user ID", err)
		return
	}

	_, err = apiCfg.db.GetUserByID(context.Background(), id)
	if err == sql.ErrNoRows {
		respondWithError(w, 404, "Couldn't get user by ID; empty rows", err)
		return
	}
	if err != nil {
		respondWithError(w, 500, "Couldn't get user by ID, other error:", err)
		return
	}

	err = apiCfg.db.UpgradeUserByID(context.Background(), id)
	if err != nil {
		respondWithError(w, 404, "Couldn't upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
