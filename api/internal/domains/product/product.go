package product

import (
	"github.com/Corray333/bumblebee/internal/domains/product/external"
	"github.com/Corray333/bumblebee/internal/domains/product/repository"
	"github.com/Corray333/bumblebee/internal/domains/product/service"
	"github.com/Corray333/bumblebee/internal/domains/product/transport"
	"github.com/Corray333/bumblebee/internal/storage"
	"github.com/Corray333/bumblebee/internal/telegram"
	"github.com/go-chi/chi/v5"
)

type DomainController struct {
	repo      *repository.DomainRepository
	service   *service.ProductService
	transport *transport.ProductTransport
	tg        *telegram.TelegramClient
}

func NewDomainController(router *chi.Mux, store *storage.Storage, tg *telegram.TelegramClient) *DomainController {
	repo := repository.New(store)

	external := external.New(tg.Bot)

	service := service.New(repo, external)

	transport := transport.New(router, service, tg.Bot)

	return &DomainController{
		repo:      repo,
		service:   service,
		transport: transport,
		tg:        tg,
	}
}

func (c *DomainController) Build() {
	c.transport.RegisterRoutes()
	c.service.UpdateProducts()
	go c.transport.RunTelegram()
}

func (c *DomainController) Run() {
	c.service.Run()
}
