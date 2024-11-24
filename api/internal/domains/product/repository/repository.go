package repository

import (
	"context"
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
			INSERT INTO products (name, weight, lifetime)
			VALUES ($1, $2, $3)
			ON CONFLICT (product_id) DO UPDATE SET weight = $2, lifetime = $3
		`, product.Name, product.Weight, product.Lifetime)
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

func (r *DomainRepository) GetProducts(ctx context.Context, offset int) (products []entities.Product, err error) {
	tx, _, err := r.getTx(ctx)
	if err != nil {
		return nil, err
	}

	err = tx.SelectContext(ctx, &products, `
		SELECT *
		FROM products OFFSET $1
	`, offset)
	if err != nil {
		slog.Error("failed to get products: " + err.Error())
		return nil, err
	}

	return products, nil
}

func (r *DomainRepository) CreateOrder(ctx context.Context, order *entities.Order) (err error) {
	tx, isNew, err := r.getTx(ctx)
	if err != nil {
		return err
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
		return err
	}

	err = row.Scan(&order.ID)

	for _, item := range order.Products {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO order_items (order_id, product_id, amount)
			VALUES ($1, $2, $3)
		`, order.ID, item.ID, item.Amount)
		if err != nil {
			slog.Error("failed to insert order item: " + err.Error())
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
