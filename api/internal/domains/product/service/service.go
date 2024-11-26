package service

import (
	"context"
	"os"
	"strconv"

	"github.com/Corray333/bumblebee/internal/domains/product/entities"
)

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	SetProducts(ctx context.Context, products []entities.Product) error
	CreateProduct(ctx context.Context, product *entities.Product) (id int64, err error)
	EditProduct(ctx context.Context, product *entities.Product) (err error)
	DeleteProduct(ctx context.Context, productID int64) (err error)
	GetProducts(ctx context.Context, offset int) (products []entities.Product, err error)
	ReorderProduct(ctx context.Context, productID int64, newPosition int) (err error)

	CreateOrder(ctx context.Context, order *entities.Order) (orderID int64, err error)

	GetProductsInOrder(ctx context.Context, order *entities.Order) (*entities.Order, error)

	GetManagerByPhoneOrEmail(ctx context.Context, manager *entities.Manager) (*entities.Manager, error)
	GetManagerByID(ctx context.Context, manager *entities.Manager) (*entities.Manager, error)
	SetManager(ctx context.Context, manager *entities.Manager) (err error)
}

type external interface {
	GetProducts() (products []entities.Product, err error)
	SendNewOrderMessage(ctx context.Context, order *entities.Order) error
}

type fileManager interface {
	SaveImage(file []byte, name string) error
	UploadImage(file []byte, name string) (string, error)
}

type ProductService struct {
	repo        repository
	external    external
	fileManager fileManager
}

func New(repo repository, external external, fileManager fileManager) *ProductService {
	s := &ProductService{
		repo:        repo,
		external:    external,
		fileManager: fileManager,
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
	products, err = s.repo.GetProducts(context.Background(), offset)
	if err != nil {
		return nil, err
	}

	for i := range products {
		products[i].Img = os.Getenv("BASE_URL") + products[i].Img
	}

	return products, nil
}

func (s *ProductService) PlaceOrder(order *entities.Order) error {
	orderID, err := s.repo.CreateOrder(context.Background(), order)
	if err != nil {
		return err
	}

	order.ID = orderID

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

func (s *ProductService) ReorderProduct(ctx context.Context, productID int64, newPosition int) error {
	return s.repo.ReorderProduct(ctx, productID, newPosition)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *entities.Product, photo []byte) error {
	img, err := s.fileManager.UploadImage(photo, "")
	if err != nil {
		return err
	}

	product.Img = img[2:]

	_, err = s.repo.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) EditProduct(ctx context.Context, product *entities.Product, photo []byte) error {
	if len(photo) > 0 {
		img, err := s.fileManager.UploadImage(photo, strconv.Itoa(int(product.ID)))
		if err != nil {
			return err
		}

		product.Img = img[2:]
	}

	err := s.repo.EditProduct(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, productID int64) error {
	return s.repo.DeleteProduct(ctx, productID)
}
