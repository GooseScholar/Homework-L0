package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

//const dsn = "host=localhost port=5432 user=homework password=homework1 dbname=postgres sslmode=disable"

const dsn = "host=localhost port=5432 user=postgres password=qwerty1 dbname=postgres sslmode=disable"

func NewDB(ctx context.Context) (*DB, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{
		pool: pool,
	}, nil
}
