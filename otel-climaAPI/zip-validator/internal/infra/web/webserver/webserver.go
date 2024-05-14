package webserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type WebServer struct {
	Router        chi.Router
	Handlers      []HandlerFunc
	WebServerPort string
	Server        *http.Server
}

type HandlerFunc struct {
	Method  string
	Handler http.HandlerFunc
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) MountMiddlewares() {
	// Middlewares
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.AllowContentType("application/json"))
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

func (s *WebServer) Start() error {
	fmt.Println("Starting web server on port", s.WebServerPort)

	// Configure a server instance
	s.Server = &http.Server{
		Addr:    s.WebServerPort,
		Handler: s.Router,
	}

	// Start the server in a separate goroutine
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("could not listen on %s: %v\n", s.WebServerPort, err)
		}
	}()

	// Listen for shutdown signals
	s.listenForShutdown()

	return nil
}

func (s *WebServer) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down web server...")

	// Create a context with a timeout to allow for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Shutdown the server
	return s.Server.Shutdown(shutdownCtx)
}

func (s *WebServer) listenForShutdown() {
	// Create a channel to listen for OS signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Block until a signal is received
	<-signals

	// Create a context to initiate shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initiate shutdown
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("Failed to gracefully shutdown the server:", err)
	}

	fmt.Println("Server has been gracefully shutdown")
}
