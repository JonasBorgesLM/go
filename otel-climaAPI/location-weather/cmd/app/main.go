package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/configs"
	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/internal/infra/web/controllers"
	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/internal/infra/web/handlers"
	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/internal/infra/web/webserver"
	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/internal/service"
)

func main() {
	// Carrega a configuração do ambiente
	if err := configs.LoadConfig("."); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Cria um canal para receber sinais de interrupção
	signChannel := make(chan os.Signal, 1)
	signal.Notify(signChannel, os.Interrupt)

	// Cria um contexto para monitorar sinais de interrupção
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Configura e inicializa o provedor de rastreamento OpenTelemetry
	shutdown, err := configs.Setup()
	if err != nil {
		log.Fatalf("failed to setup tracer provider: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown tracer provider: ", err)
		}
	}()

	// Cria uma instância do servidor web
	server := webserver.NewWebServer(":8081")
	server.MountMiddlewares()

	viaCepService := service.NewViaCepService()
	weatherApiService := service.NewWeatherApiService()

	// Cria uma instância do manipulador HTTP
	handler := handlers.NewHandler(viaCepService, weatherApiService)

	// Cria uma instância do controlador e define as rotas
	controller := controllers.NewController(server.Router, handler)
	controller.AddRoutes()

	// Inicia o servidor em uma gorrotina separada
	go server.Start()

	// Aguarda os sinais de interrupção ou cancelamento do contexto
	select {
	case <-signChannel:
		log.Println("shutting down server gracefully...")
	case <-ctx.Done():
		log.Println("shutting down server...")
	}

	// Encerra o servidor e aguarda a conclusão
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server shutdown complete")
}
