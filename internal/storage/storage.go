package storage

import (
	"HappyKod/service-api-gofermart/internal/models"
	"context"
)

type GoferStorage interface {
	Ping() error
	Close() error
	AddUser(ctx context.Context, user models.User) error
	AuthenticationUser(ctx context.Context, user models.User) (bool, error)
	GetOrder(ctx context.Context, numberOrder string) (models.Order, error)
	GetManyOrders(ctx context.Context, userLogin string) ([]models.Order, error)
	AddOrder(ctx context.Context, numberOrder string, order models.Order) error
	GetOrdersByProcess(ctx context.Context) ([]models.Order, error)
	UpdateOrder(ctx context.Context, loyaltyPoint models.LoyaltyPoint) error
	GetUserBalance(ctx context.Context, userLogin string) (float64, float64, error)
	AddWithdraw(ctx context.Context, withdraw models.Withdraw) error
	GetManyWithdraws(ctx context.Context, userLogin string) ([]models.Withdraw, error)
}
