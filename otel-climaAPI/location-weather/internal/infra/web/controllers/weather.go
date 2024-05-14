package controllers

import (
	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/internal/infra/web/handlers"
	"github.com/go-chi/chi/v5"
)

// Controller representa um controlador para lidar com solicitações HTTP.
type Controller struct {
	Router  chi.Router
	Handler *handlers.Handler
}

// NewController cria uma nova instância de Controller.
func NewController(router chi.Router, handler *handlers.Handler) *Controller {
	return &Controller{
		Router:  router,
		Handler: handler,
	}
}

// AddRoutes adiciona rotas ao controlador.
func (c *Controller) AddRoutes() {
	c.Router.Route("/", func(r chi.Router) {
		r.Get("/", c.Handler.GetTemperatures)
	})
}
