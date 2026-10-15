package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	var server *http.Server

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz/", healthHandler)

	server = &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}

func healthHandler(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.WriteHeader(200)
	_, err := resWriter.Write([]byte("OK"))
	if err != nil {
		log.Fatal("could not write response message")
	}

	resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
}
