package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

const dsn = "host=localhost port=5432 user=homework password=homework1 dbname=postgres sslmode=disable"

//Подключение к базе данных от имени пользователя homework
func NewDB(ctx context.Context) (*DB, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{
		pool: pool,
	}, nil
}
