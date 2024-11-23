package transport

import (
	"github.com/go-chi/chi/v5"
)

type DomainTransport struct {
	router  *chi.Mux
	service service
}

type service interface {
}

func New(router *chi.Mux, service service) *DomainTransport {
	return &DomainTransport{
		router:  router,
		service: service,
	}
}

func (t *DomainTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {

	})
}
