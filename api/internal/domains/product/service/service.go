package service

import (
	"context"

	"github.com/Corray333/bumblebee/internal/domains/product/entities"
)

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	SetProducts(ctx context.Context, products []entities.Product) error
	GetProducts(ctx context.Context, offset int) (products []entities.Product, err error)
	CreateOrder(ctx context.Context, order *entities.Order) error

	GetProductsInOrder(ctx context.Context, order *entities.Order) (*entities.Order, error)

	GetManagerByPhoneOrEmail(ctx context.Context, manager *entities.Manager) (*entities.Manager, error)
	GetManagerByID(ctx context.Context, manager *entities.Manager) (*entities.Manager, error)
	SetManager(ctx context.Context, manager *entities.Manager) (err error)
}

type external interface {
	GetProducts() (products []entities.Product, err error)
	SendNewOrderMessage(ctx context.Context, order *entities.Order) error
}

type ProductService struct {
	repo     repository
	external external
}

func New(repo repository, external external) *ProductService {
	s := &ProductService{
		repo:     repo,
		external: external,
	}
	return s
}

func (s *ProductService) Run() {
}

func (s *ProductService) UpdateProducts() error {
	products, err := s.external.GetProducts()
	if err != nil {
		return err
	}

	if err := s.repo.SetProducts(context.Background(), products); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) GetProducts(offset int) (products []entities.Product, err error) {
	return s.repo.GetProducts(context.Background(), offset)
}

func (s *ProductService) PlaceOrder(order *entities.Order) error {
	err := s.repo.CreateOrder(context.Background(), order)
	if err != nil {
		return err
	}

	manager, err := s.repo.GetManagerByPhoneOrEmail(context.Background(), &order.Manager)
	if err != nil {
		return err
	}

	order.Manager = *manager

	order, err = s.repo.GetProductsInOrder(context.Background(), order)
	if err != nil {
		return err
	}

	return s.external.SendNewOrderMessage(context.Background(), order)
}

func (s *ProductService) GetManagerByPhoneOrEmail(ctx context.Context, manager *entities.Manager) (*entities.Manager, error) {
	return s.repo.GetManagerByPhoneOrEmail(ctx, manager)
}

func (s *ProductService) GetManagerByID(ctx context.Context, manager *entities.Manager) (*entities.Manager, error) {
	return s.repo.GetManagerByID(ctx, manager)
}

func (s *ProductService) SetManager(ctx context.Context, manager *entities.Manager) error {
	return s.repo.SetManager(ctx, manager)
}
