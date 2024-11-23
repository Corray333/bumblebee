package storage

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TxKey struct{}

// Errors
var (
	ErrInvalidTxType = fmt.Errorf("invalid transaction type")
)

type Storage struct {
	db *sqlx.DB
}

func (s *Storage) DB() *sqlx.DB {
	return s.db
}

func New() (*Storage, error) {
	dbName := os.Getenv("POSTGRES_DB_NAME")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbHost := os.Getenv("POSTGRES_HOST")

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}
