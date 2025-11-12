package main

import (
	"net/http"
	"time"

	"github.com/jifpbj/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not able to parse Authorization Header", err)
		return
	}

	tokenData, err := apiCfg.db.RetrieveRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	TokenExpirationDuration := time.Duration(3600) * time.Second

	token, err := auth.MakeJWT(tokenData.UserID, apiCfg.secret, TokenExpirationDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}

	type payload struct {
		Token string `json:"token"`
	}

	out := payload{
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, out)
}
