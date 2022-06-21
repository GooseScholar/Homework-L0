package service

import (
	"context"
	"homework-l0/internal/models"
)

type OrderService struct {
	repo Repository
}

func NewOrderService(repo Repository) *OrderService {
	return &OrderService{repo: repo}
}

//Запись данных в бд
func (os *OrderService) PutOrder(ctx context.Context, ord *models.Orders) (err error) {

	return
}

//Получение данных из бд
func (os *OrderService) GetOrder(ctx context.Context, Order_uid string) (ord *models.Orders, err error) {

	return
}
