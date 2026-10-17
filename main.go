package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	mux.HandleFunc("POST /api/validate_chirp", chirpHandler)

	server = &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}

func chirpHandler(resWriter http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Error string `json:"error"`
		Body  string `json:"body"`
	}

	resWriter.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)

	var theChirp chirp

	if err := decoder.Decode(&theChirp); err != nil {
		resWriter.WriteHeader(500)
		responseBody := returnVals{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(responseBody)
		if err != nil {
			log.Fatal("could not write to json")
		}

		_, err = resWriter.Write(dat)
		if err != nil {
			log.Fatal("cant write request")
		}
		return
	}

	if len(theChirp.Body) > 140 {
		resWriter.WriteHeader(400)
		responseBody := returnVals{
			Error: "chirp is tooo long",
		}
		dat, err := json.Marshal(responseBody)
		if err != nil {
			log.Fatal("could not write to json")
		}
		_, err = resWriter.Write(dat)
		if err != nil {
			log.Fatal("cant write request")
		}
		return
	}

	theChirp.Body = dirtyDudesMustDie(theChirp.Body)

	resWriter.WriteHeader(200)
	responseBody := returnVals{
		Body: theChirp.Body,
	}
	dat, err := json.Marshal(responseBody)
	if err != nil {
		log.Fatal("could not write to json")
	}
	_, err = resWriter.Write(dat)
	if err != nil {
		log.Fatal("cant write request")
	}
}

func dirtyDudesMustDie(body string) string {
	profane := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	var finalBody strings.Builder

	words := strings.Fields(body)

	for _, word := range words {
		if profane[strings.ToLower(string(word))] {
			finalBody.WriteString("****")
		} else {
			finalBody.WriteString(string(word))
		}
		finalBody.WriteRune(' ')
	}
	return finalBody.String()
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
