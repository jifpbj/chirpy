package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jifpbj/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("author_id")
	var chirps []database.Chirp

	if str != "" {
		id, err := uuid.Parse(str)
		if err != nil {
			respondWithError(w, 400, "invalid author id", err)
			return
		}
		cs, err := apiCfg.db.RetrieveChirpsByUserID(r.Context(), id)
		if err != nil {
			respondWithError(w, 500, "err retrieving chirps by author id", err)
			return
		}
		chirps = cs
	} else {
		cs, err := apiCfg.db.RetrieveChirps(context.Background())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
		chirps = cs
	}

	output := make([]Chirp, 0)

	for _, chirp := range chirps {
		out := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		output = append(output, out)
	}
	respondWithJSON(w, 200, output)
}

func (apiCfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid ID", err)
		return
	}

	chirp, err := apiCfg.db.RetrieveChirpByID(context.Background(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Chirp with ID %s not found", chirpID), err)
		}
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp by ID", err)
		return
	}
	out := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, out)
}
