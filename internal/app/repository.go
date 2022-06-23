package app

import (
	"context"
	"homework-l0/internal/models"
)

type Repository interface {
	PutOrder(context.Context, *models.Orders) error
}
