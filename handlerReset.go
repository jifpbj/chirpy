package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if apiCfg.platform != "dev" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := apiCfg.db.ResetUser(req.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset users", err)
		return
	}

	apiCfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hits reset back to 0"))
	if err != nil {
		fmt.Printf("error writing statusText: %s", err)
	}
}
