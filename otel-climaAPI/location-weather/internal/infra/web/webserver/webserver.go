package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// WebServer representa um servidor web.
type WebServer struct {
	Router        chi.Router
	Handlers      []HandlerFunc
	WebServerPort string // Porta do servidor web
	server        *http.Server
}

// HandlerFunc representa um manipulador de rota.
type HandlerFunc struct {
	Method  string
	Handler http.HandlerFunc
}

// NewWebServer cria uma nova instância do servidor web.
func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		WebServerPort: serverPort, // Define a porta do servidor
	}
}

// MountMiddlewares adiciona middlewares ao roteador.
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
		MaxAge:           300, // 5 minutes
	}))
}

// Start inicia o servidor web.
func (s *WebServer) Start() {
	fmt.Println("Starting web server on port", s.WebServerPort)
	s.server = &http.Server{
		Addr:    s.WebServerPort,
		Handler: s.Router,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("failed to start server: %v\n", err)
		}
	}()
}

// Shutdown encerra o servidor web.
func (s *WebServer) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down web server...")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	fmt.Println("Web server shutdown complete")
	return nil
}
