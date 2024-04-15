package main

import (
	"log"
	"net/http"

	"github.com/JonasBorgesLM/go/ratelimiter/configs"
	"github.com/JonasBorgesLM/go/ratelimiter/internal/infra/webserver"
)

func main() {
	// Load configurations
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		log.Fatalf("error loading configurations")

		panic(err)
	}

	webserver := webserver.NewWebServer(configs)
	webserver.AddHandler("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server listening on %v\n", configs.ServerPort)
	webserver.Start()
}
