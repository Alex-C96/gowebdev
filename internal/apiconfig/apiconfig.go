package apiconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ApiConfig struct {
	fileserverHits int
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

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type failureDetails struct {
		Body string `json:"body"`
	}

	type successDetails struct {
		Valid bool `json:"valid"`
	}

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	chirpResp := chirp{}
	err := decoder.Decode(&chirpResp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respBody := failureDetails{
			Body: "Something went wrong",
		}
		dat, _ := json.Marshal(respBody)
		w.WriteHeader(500)
		w.Write(dat)
		return
	}
	length := len(chirpResp.Body)

	if length > 140 {
		respBody := failureDetails{
			Body: "Chirp is too long",
		}
		dat, _ := json.Marshal(respBody)
		w.WriteHeader(400)
		w.Write(dat)
	} else {
		respBody := successDetails{
			Valid: true,
		}
		dat, _ := json.Marshal(respBody)
		w.WriteHeader(200)
		w.Write(dat)
	}

}
