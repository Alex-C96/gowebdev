package apiconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alex-c96/gowebdev/internal/database"
	"github.com/go-chi/chi/v5"
)

type ApiConfig struct {
	fileserverHits int
	Database       *database.DB
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) GetHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf8")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`
	<html>

		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>

	</html>

	`, cfg.fileserverHits)
	w.Write([]byte(html))
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *ApiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Database.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not read chirps from DB")
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *ApiConfig) GetChirpId(w http.ResponseWriter, r *http.Request) {
	chirpId := chi.URLParam(r, "chirpID")
	num, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, 400, "invalid chirp id")
		return
	}
	chirp, err := cfg.Database.GetChirpByID(num)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}
	respondWithJSON(w, 200, chirp)
}

func (cfg *ApiConfig) PostChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type successDetails struct {
		Body string `json:"body"`
		ID   int    `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpResp := chirp{}
	err := decoder.Decode(&chirpResp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters")
		return
	}
	length := len(chirpResp.Body)
	const maxChirpLength = 140

	if length > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	newChirp := successDetails{
		Body: chirpResp.Body,
		ID:   cfg.Database.ChirpId,
	}

	respondWithJSON(w, 201, newChirp)

	responseChirp, err := cfg.Database.CreateChirp(chirpResp.Body)
	if err != nil {
		log.Printf("Db input failure %v", err)
	}
	fmt.Println(responseChirp)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
