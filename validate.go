package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleaned := wordReplacement(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleaned})
}

func wordReplacement(txt string) string {
	words := strings.Split(txt, " ")
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	for i, word := range words {
		lower := strings.ToLower(word)
		if slices.Contains(badWords, lower) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
