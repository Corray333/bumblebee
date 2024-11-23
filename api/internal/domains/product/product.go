package product

import (
	"github.com/Corray333/bumblebee/internal/domains/product/repository"
	"github.com/Corray333/bumblebee/internal/domains/product/service"
	"github.com/Corray333/bumblebee/internal/domains/product/transport"
	"github.com/Corray333/bumblebee/internal/storage"
	"github.com/go-chi/chi/v5"
)

type DomainController struct {
	repo      repository.DomainRepository
	service   service.DomainService
	transport transport.DomainTransport
}

func NewDomainController(router *chi.Mux, store *storage.Storage) *DomainController {
	repo := repository.New(store)
	service := service.New(repo)
	transport := transport.New(router, service)

	return &DomainController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *DomainController) Build() {
	c.transport.RegisterRoutes()
}

func (c *DomainController) Run() {
	c.service.Run()
}
