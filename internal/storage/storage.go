package storage

import "HappyKod/service-api-gofermart/internal/models"

type GoferStorage interface {
	Ping() error
	Close() error
	UniqLoginUser(login string) (bool, error)
	AddUser(user models.User) error
	AuthenticationUser(user models.User) (bool, error)
	GetOrder(numberOrder string) (models.Order, error)
	GetManyOrders(userLogin string) ([]models.Order, error)
	AddOrder(numberOrder string, order models.Order) error
	GetOrdersByProcess() ([]models.Order, error)
	UpdateOrder(loyaltyPoint models.LoyaltyPoint) error
	GetUserBalance(userLogin string) (float64, float64, error)
	AddWithdraw(withdraw models.Withdraw) error
}
