package main

import (
	"log"
	"net/http"

	"github.com/alex-c96/gowebdev/internal/apiconfig"
	"github.com/go-chi/chi/v5"
)

func main() {
	const port = "8080"
	const filePathRoot = "."
	r := chi.NewRouter()
	apiR := chi.NewRouter()
	adminR := chi.NewRouter()
	apiCfg := new(apiconfig.ApiConfig)
	fsHandler := apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	r.Handle("/app/", fsHandler)
	r.Handle("/app*", fsHandler)
	apiR.Get("/healthz", handlerReadiness)
	apiR.Get("/reset", apiCfg.Reset)
	adminR.Get("/metrics", apiCfg.GetHits)
	r.Mount("/admin", adminR)
	r.Mount("/api", apiR)
	corsr := middlewareLogger(middlewareCors(r))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsr,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
