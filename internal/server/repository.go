package server

import (
	"context"
	"homework-l0/internal/cache"
	"homework-l0/internal/models"
)

type Repository interface {
	GetOrder(context.Context, string) (*models.Orders, error)
	GetInitialCache(context.Context) (*cache.Cache, error)
}
