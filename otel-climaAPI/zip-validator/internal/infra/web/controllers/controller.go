package controllers

import (
	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/internal/infra/web/handlers"
	"github.com/go-chi/chi/v5"
)

// Controller representa o controlador responsável por configurar as rotas e associá-las aos manipuladores
type Controller struct {
	router  chi.Router
	handler *handlers.Handler
}

// NewController cria uma nova instância do Controller
func NewController(router chi.Router, handler *handlers.Handler) *Controller {
	return &Controller{
		router:  router,
		handler: handler,
	}
}

// Route configura as rotas do controlador
func (c *Controller) Route() {
	c.router.Route("/", func(r chi.Router) {
		r.Post("/", c.handler.GetTemperatures)
	})
}
