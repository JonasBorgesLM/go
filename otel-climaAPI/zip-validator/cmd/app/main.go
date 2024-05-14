package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	configs "github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/config"
	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/internal/infra/web/controllers"
	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/internal/infra/web/handlers"
	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/internal/infra/web/webserver"
	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/internal/service"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configura e inicializa o provedor de rastreamento OpenTelemetry
	shutdown, err := configs.Setup()
	if err != nil {
		log.Fatalf("failed to setup tracer provider: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer provider: %v", err)
		}
	}()

	// Cria uma instância do servidor web
	server := webserver.NewWebServer(":8080")
	server.MountMiddlewares()

	// Cria uma instância do serviço de obtenção de temperatura
	apiService := service.NewGetTemperatureService()

	// Cria uma instância do manipulador HTTP
	handler := handlers.NewHandler(apiService)

	// Cria uma instância do controlador e define as rotas
	controller := controllers.NewController(server.Router, handler)
	controller.Route()

	// Inicia o servidor em uma gorrotina separada
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Aguarda os sinais de interrupção ou cancelamento do contexto
	select {
	case <-signChannel:
		log.Println("received interrupt signal, shutting down server...")
	case <-ctx.Done():
		log.Println("context cancelled, shutting down server...")
	}

	// Encerra o servidor e aguarda a conclusão
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server shutdown complete")
}
