package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Corray333/bumblebee/internal/domains/product/entities"
	"github.com/Corray333/bumblebee/internal/storage"
	"github.com/jmoiron/sqlx"
)

type DomainRepository struct {
	db *sqlx.DB
}

func New(store *storage.Storage) *DomainRepository {
	return &DomainRepository{
		db: store.DB(),
	}
}

func (r *DomainRepository) Begin(ctx context.Context) (context.Context, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, storage.TxKey{}, tx), nil
}

func (r *DomainRepository) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(storage.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Commit()
}

func (r *DomainRepository) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(storage.TxKey{}).(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx.Rollback()
}

func (r *DomainRepository) getTx(ctx context.Context) (tx *sqlx.Tx, isNew bool, err error) {
	txRaw := ctx.Value(storage.TxKey{})
	if txRaw != nil {
		var ok bool
		tx, ok = txRaw.(*sqlx.Tx)
		if !ok {
			slog.Error("invalid transaction type")
			return nil, false, storage.ErrInvalidTxType
		}
	}
	if tx == nil {
		tx, err = r.db.BeginTxx(ctx, nil)
		if err != nil {
			slog.Error("failed to begin transaction: " + err.Error())
			return nil, false, err
		}

		return tx, true, nil
	}

	return tx, false, nil
}

func (r *DomainRepository) SetProducts(ctx context.Context, products []entities.Product) error {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM products")
	if err != nil {
		slog.Error("failed to clean products: " + err.Error())
		return err
	}

	for _, product := range products {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO products (description, position)
			VALUES ($1, (SELECT COALESCE(MAX(position), 0) + 1 FROM products))
			ON CONFLICT (product_id) DO UPDATE SET description = $1
		`, product.Description)
		if err != nil {
			slog.Error("failed to insert product: " + err.Error())
			return err
		}
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}

	return nil
}

func (r *DomainRepository) CreateProduct(ctx context.Context, product *entities.Product) (id int64, err error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return 0, err
	}
	if isNew {
		defer tx.Rollback()
	}

	row := tx.QueryRow(`
		INSERT INTO products (description, img, position)
		VALUES ($1, $2, (SELECT COALESCE(MAX(position), 0) + 1 FROM products))
		ON CONFLICT (product_id) DO UPDATE SET description = $1 RETURNING product_id
	`, product.Description, product.Img)

	err = row.Scan(&id)
	if err != nil {
		slog.Error("failed to insert product: " + err.Error())
		return 0, err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return 0, err
		}
	}

	return id, nil
}

func (r *DomainRepository) ReorderProduct(ctx context.Context, productID int64, newPosition int) (err error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	// Получаем текущую позицию продукта
	var currentPosition int
	err = tx.QueryRow(`SELECT position FROM products WHERE product_id = $1`, productID).Scan(&currentPosition)
	if err != nil {
		slog.Error("failed to get current position: " + err.Error())
		return err
	}

	// Обновляем позиции продуктов
	if currentPosition < newPosition {
		_, err = tx.Exec(`
            UPDATE products
            SET position = position - 1
            WHERE position > $1 AND position <= $2
        `, currentPosition, newPosition)
	} else {
		_, err = tx.Exec(`
            UPDATE products
            SET position = position + 1
            WHERE position >= $1 AND position < $2
        `, newPosition, currentPosition)
	}
	if err != nil {
		slog.Error("failed to update positions: " + err.Error())
		return err
	}

	// Устанавливаем новую позицию для продукта
	_, err = tx.Exec(`
        UPDATE products
        SET position = $1
        WHERE product_id = $2
    `, newPosition, productID)
	if err != nil {
		slog.Error("failed to set new position: " + err.Error())
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}

	return nil
}

func (r *DomainRepository) GetProducts(ctx context.Context, offset int) (products []entities.Product, err error) {
	err = r.db.Select(&products, `
		SELECT *
		FROM products ORDER BY position OFFSET $1
	`, offset)
	if err != nil {
		slog.Error("failed to get products: " + err.Error())
		return nil, err
	}

	return products, nil
}

func (r *DomainRepository) CreateOrder(ctx context.Context, order *entities.Order) (orderID int64, err error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return 0, err
	}
	if isNew {
		defer tx.Rollback()
	}

	row := tx.QueryRowx(`
		INSERT INTO orders (date, customer_phone, customer_name, customer_address)
		VALUES ($1, $2, $3, $4)
		RETURNING order_id
	`, order.Date, order.Customer.Phone, order.Customer.Name, order.Customer.Address)
	if err != nil {
		slog.Error("failed to insert order: " + err.Error())
		return 0, err
	}

	err = row.Scan(&order.ID)

	for _, item := range order.Products {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO order_product (order_id, product_id, amount)
			VALUES ($1, $2, $3)
		`, order.ID, item.ID, item.Amount)
		if err != nil {
			slog.Error("failed to insert order item: " + err.Error())
			return 0, err
		}
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return 0, err
		}
	}

	return order.ID, nil
}

func (r *DomainRepository) SetManager(ctx context.Context, manager *entities.Manager) (err error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO managers (manager_id, state, phone, email)
		VALUES ($1, $2, $3, $4) ON CONFLICT (manager_id) DO UPDATE SET state = $2, phone = $3, email = $4
	`, manager.ID, manager.State, manager.Phone, manager.Email)
	if err != nil {
		slog.Error("failed to insert manager: " + err.Error())
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}

	return nil
}

func (r *DomainRepository) GetManagerByPhoneOrEmail(ctx context.Context, manager *entities.Manager) (*entities.Manager, error) {

	fmt.Printf("manager: %+v\n", manager)
	err := r.db.Get(manager, `
		SELECT *
		FROM managers
		WHERE phone = $1 OR email = $2
	`, manager.Phone, manager.Email)
	if err != nil {
		slog.Error("failed to get manager: " + err.Error())
		return nil, err
	}

	return manager, nil
}

func (r *DomainRepository) GetManagerByID(ctx context.Context, manager *entities.Manager) (*entities.Manager, error) {

	err := r.db.Get(manager, `
		SELECT *
		FROM managers
		WHERE manager_id = $1
	`, manager.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to get manager: " + err.Error())
		return nil, err
	}

	return manager, nil
}

func (r *DomainRepository) GetProductsInOrder(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	for i := range order.Products {
		product := entities.Product{}
		err := r.db.Get(&product, `
			SELECT *
			FROM products
			WHERE product_id = $1
		`, order.Products[i].ID)
		if err != nil {
			slog.Error("failed to get product: " + err.Error())
			return nil, err
		}
		order.ProductList = append(order.ProductList, product)
	}

	return order, nil
}

func (r *DomainRepository) EditProduct(ctx context.Context, product *entities.Product) (err error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	fmt.Printf("product: %+v", product)

	_, err = tx.ExecContext(ctx, `
		UPDATE products
		SET description = $1, img = $2
		WHERE product_id = $3
	`, product.Description, product.Img, product.ID)
	if err != nil {
		slog.Error("failed to update product: " + err.Error())
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}

	return nil
}

func (r *DomainRepository) DeleteProduct(ctx context.Context, productID int64) (err error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
	}
	if isNew {
		defer tx.Rollback()
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM products
		WHERE product_id = $1
	`, productID)
	if err != nil {
		slog.Error("failed to delete product: " + err.Error())
		return err
	}

	if isNew {
		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction: " + err.Error())
			return err
		}
	}

	return nil
}
