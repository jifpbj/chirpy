package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	apiCfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hits reset back to 0"))
	if err != nil {
		fmt.Printf("error writing statusText: %s", err)
	}
}
