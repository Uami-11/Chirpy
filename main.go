package main

import (
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()

	serveMux.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("could not start up server")
	}
}
