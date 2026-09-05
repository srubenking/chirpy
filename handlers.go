package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = atomic.Int32{}
	w.Write([]byte("metrics reset"))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK")) // handle errors here later?
}

func validate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	body := ReqBody{}

	if err := decoder.Decode(&body); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not decode POST JSON: %s", err))
		return
	}

	if chirpLength := len(body.Body); chirpLength > 140 {
		respondWithError(w, 400, fmt.Sprintf("Chirp is too long (%d characters)", chirpLength))
		return
	}

	cleanBody := CleanBody{
		CleanedBody: cleanChirp(body.Body),
	}
	respondWithJSON(w, 200, cleanBody)
}
