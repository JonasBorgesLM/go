package main

import (
	"log"
	"net/http"

	"github.com/JonasBorgesLM/go/clound-run/internal/infra/webserver"
)

func main() {
	port := "8080"

	webserver := webserver.NewWebServer(port)
	webserver.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server listening on %v\n", port)
	webserver.Start()
}
