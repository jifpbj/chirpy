package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/jifpbj/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	var chirps []database.Chirp

	if authorID != "" {
		id, err := uuid.Parse(authorID)
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

	// sort Order based on sortOrder
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		respondWithError(w, 400, "invalid sort order", errors.New("sort must be 'asc' or 'desc'"))
		return
	}
	if sortOrder == "" || sortOrder == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	} else if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
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
