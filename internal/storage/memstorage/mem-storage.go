package memstorage

import (
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"errors"
	"github.com/google/uuid"
	"sync"
)

type MemStorage struct {
	userCash  map[uuid.UUID]models.User
	orderCash map[string]models.Order
	mu        *sync.RWMutex
}

func New() (*MemStorage, error) {
	return &MemStorage{
		userCash:  make(map[uuid.UUID]models.User),
		orderCash: make(map[string]models.Order),
		mu:        new(sync.RWMutex),
	}, nil
}

func (MS *MemStorage) Ping() error {
	return nil
}

func (MS *MemStorage) Close() error {
	return nil
}

func (MS *MemStorage) AddUser(user models.User) error {
	MS.mu.Lock()
	defer MS.mu.Unlock()
	MS.userCash[uuid.New()] = user
	return nil
}

func (MS *MemStorage) UniqLoginUser(login string) (bool, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	for _, v := range MS.userCash {
		if v.Login == login {
			return false, nil
		}
	}
	return true, nil
}

func (MS *MemStorage) AuthenticationUser(user models.User) (bool, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	for _, v := range MS.userCash {
		if v.Login == user.Login && v.Password == user.Password {
			return true, nil
		}
	}
	return false, nil
}

func (MS *MemStorage) GetOrder(numberOrder string) (models.Order, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	return MS.orderCash[numberOrder], nil
}

func (MS *MemStorage) GetManyOrders(userLogin string) ([]models.Order, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	var orders []models.Order
	for _, v := range MS.orderCash {
		if v.UserLogin == userLogin {
			orders = append(orders, v)
		}
	}
	return orders, nil
}

func (MS *MemStorage) AddOrder(numberOrder string, order models.Order) error {
	MS.mu.Lock()
	defer MS.mu.Unlock()
	if MS.orderCash[numberOrder].NumberOrder == numberOrder {
		return errors.New(constans.ErrorNoUNIQUE)
	}
	MS.orderCash[numberOrder] = order
	return nil
}
