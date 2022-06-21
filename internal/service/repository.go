package service

import (
	"context"
	"homework-l0/internal/models"
)

type Repository interface {
	PutOrder(context.Context, *models.Orders) error
	GetOrder(context.Context, string) (*models.Orders, error)
}
