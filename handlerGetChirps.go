package main

import (
	"context"
	"net/http"
)

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := apiCfg.db.RetrieveChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
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
