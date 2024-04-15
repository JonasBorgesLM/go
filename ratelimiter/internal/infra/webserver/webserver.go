package webserver

import (
	"fmt"
	"net/http"

	"github.com/JonasBorgesLM/go/ratelimiter/configs"
	"github.com/JonasBorgesLM/go/ratelimiter/internal/middleware/ratelimiter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router   chi.Router
	Handlers map[string]http.HandlerFunc
	Configs  *configs.Conf
}

func NewWebServer(config *configs.Conf) *WebServer {
	return &WebServer{
		Router:   chi.NewRouter(),
		Handlers: make(map[string]http.HandlerFunc),
		Configs:  config,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	// Create rate limiter middleware
	rateLimiter := ratelimiter.NewRateLimiter(s.Configs)

	s.Router.Use(rateLimiter)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	serverPort := fmt.Sprint(":", s.Configs.ServerPort)
	http.ListenAndServe(serverPort, s.Router)
}
