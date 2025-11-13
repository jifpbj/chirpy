package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/jifpbj/chirpy/internal/auth"
	"github.com/jifpbj/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	BearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(BearerToken, apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid ID", err)
		return
	}

	chirp, err := apiCfg.db.RetrieveChirpByID(r.Context(), chirpUUID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to retrieve chirp", err)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you do not have permission to delete this chirp", nil)
		return
	}

	err = apiCfg.db.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
		ID:     chirpUUID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
