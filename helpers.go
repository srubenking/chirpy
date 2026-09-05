package main

import (
	"encoding/json"
	"log"
	"net/http"
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
