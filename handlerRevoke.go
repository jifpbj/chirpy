package main

import (
	"net/http"

	"github.com/jifpbj/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not able to parse Authorization Header", err)
	}
	err = apiCfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
