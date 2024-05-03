package main

import (
	"log"

	"github.com/JonasBorgesLM/go/clound-run/internal/infra/weather"
	"github.com/JonasBorgesLM/go/clound-run/internal/infra/webserver"
)

func main() {
	port := "8080"

	webserver := webserver.NewWebServer(port)
	webserver.AddHandler("/clima/{cep}", weather.ClimaHandler)

	log.Printf("Server listening on %v\n", port)
	webserver.Start()
}
