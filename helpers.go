package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type ReqBody struct {
	Body string `json:"body"`
}
type RespErr struct {
	Error string `json:"error"`
}
type RespValid struct {
	Valid bool `json:"valid"`
}
type CleanBody struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	err := RespErr{
		Error: msg,
	}
	respondWithJSON(w, code, err)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// handle marshalling errors in here, respondWithError is only for JSON responses
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func cleanChirp(msg string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	words := strings.Split(msg, " ")
	cleanWords := make([]string, 0)

	for _, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			cleanWords = append(cleanWords, "****")
		} else {
			cleanWords = append(cleanWords, word)
		}
	}
	return strings.Join(cleanWords, " ")
}
