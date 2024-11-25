package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Corray333/bumblebee/internal/domains/product/entities"
	"github.com/go-chi/chi/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProductTransport struct {
	router  *chi.Mux
	service service
	tg      *tgbotapi.BotAPI
}

type service interface {
	GetProducts(offset int) (products []entities.Product, err error)
	PlaceOrder(order *entities.Order) error
	ReorderProduct(ctx context.Context, productID int64, newPosition int) error
	CreateProduct(ctx context.Context, product *entities.Product) error
	EditProduct(ctx context.Context, product *entities.Product) error
	DeleteProduct(ctx context.Context, productID int64) error

	GetManagerByPhoneOrEmail(ctx context.Context, manager *entities.Manager) (*entities.Manager, error)
	GetManagerByID(ctx context.Context, manager *entities.Manager) (*entities.Manager, error)
	SetManager(ctx context.Context, manager *entities.Manager) error
}

func New(router *chi.Mux, service service, tg *tgbotapi.BotAPI) *ProductTransport {
	return &ProductTransport{
		router:  router,
		service: service,
		tg:      tg,
	}
}

func (t *ProductTransport) RegisterRoutes() {
	t.router.Group(func(r chi.Router) {
		r.Get("/api/products", t.getProducts)
		r.Post("/api/products", t.createProduct)
		r.Put("/api/products/{productID}", t.editProduct)
		r.Delete("/api/products/{productID}", t.deleteProduct)
		r.Put("/api/products/{productID}/reorder", t.reorderProduct)
		r.Post("/api/order", t.placeOrder)
	})
}

// @Summary      Get Products
// @Description  Retrieves a list of products from the service.
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}  entities.Product
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /api/products [get]
func (t *ProductTransport) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := t.service.GetProducts(0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary      Place Order
// @Description  Creates a new order in the system.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body      entities.Order  true  "Order Data"
// @Success      201
// @Failure      400  {string}  string  "Bad Request"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /api/order [post]
func (t *ProductTransport) placeOrder(w http.ResponseWriter, r *http.Request) {
	var order entities.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := t.service.PlaceOrder(&order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t *ProductTransport) createProduct(w http.ResponseWriter, r *http.Request) {
	var product entities.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := t.service.CreateProduct(r.Context(), &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t *ProductTransport) editProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	var product entities.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product.ID = int64(productIDInt)
	if err := t.service.EditProduct(r.Context(), &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *ProductTransport) deleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := t.service.DeleteProduct(r.Context(), int64(productIDInt)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *ProductTransport) reorderProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	newPosition := r.URL.Query().Get("new_position")
	productIDInt, err := strconv.Atoi(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newPositionInt, err := strconv.Atoi(newPosition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := t.service.ReorderProduct(r.Context(), int64(productIDInt), newPositionInt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
