package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	var server *http.Server
	var conf apiConfig

	mux.Handle("/app/", conf.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz/", healthHandler)
	mux.HandleFunc("GET /admin/metrics/", conf.metricHandler)
	mux.HandleFunc("POST /admin/reset/", conf.resetHandler)

	server = &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		handler.ServeHTTP(resWriter, req)
	})
}

func (cfg *apiConfig) resetHandler(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.WriteHeader(200)
	resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Store(0)
}

func healthHandler(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.WriteHeader(200)
	_, err := resWriter.Write([]byte("OK"))
	if err != nil {
		log.Fatal("could not write response message")
	}

	resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func (cfg *apiConfig) metricHandler(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.WriteHeader(200)
	resWriter.Header().Set("Content-Type", "text/html; charset=utf-8")

	hitString := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())

	if _, err := resWriter.Write([]byte(hitString)); err != nil {
		log.Fatal("could not write response message")
	}
}
