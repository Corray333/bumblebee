package product

import (
	"github.com/Corray333/bumblebee/internal/domains/product/external"
	"github.com/Corray333/bumblebee/internal/domains/product/repository"
	"github.com/Corray333/bumblebee/internal/domains/product/service"
	"github.com/Corray333/bumblebee/internal/domains/product/transport"
	"github.com/Corray333/bumblebee/internal/storage"
	"github.com/go-chi/chi/v5"
)

type DomainController struct {
	repo      repository.DomainRepository
	service   service.ProductService
	transport transport.DomainTransport
}

func NewDomainController(router *chi.Mux, store *storage.Storage) *DomainController {
	repo := repository.New(store)

	external := external.New()

	service := service.New(repo, external)

	transport := transport.New(router, service)

	return &DomainController{
		repo:      *repo,
		service:   *service,
		transport: *transport,
	}
}

func (c *DomainController) Build() {
	c.transport.RegisterRoutes()
	c.service.UpdateProducts()
}

func (c *DomainController) Run() {
	c.service.Run()
}
