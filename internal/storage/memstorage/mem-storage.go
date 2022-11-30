package memstorage

import (
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"context"
	"sync"

	"github.com/google/uuid"
)

type MemStorage struct {
	userCash     map[uuid.UUID]models.User
	orderCash    map[string]models.Order
	withdrawCash map[uuid.UUID]models.Withdraw
	mu           *sync.RWMutex
}

func New() (*MemStorage, error) {
	return &MemStorage{
		userCash:     make(map[uuid.UUID]models.User),
		orderCash:    make(map[string]models.Order),
		withdrawCash: make(map[uuid.UUID]models.Withdraw),
		mu:           new(sync.RWMutex),
	}, nil
}

func (MS *MemStorage) Ping() error {
	return nil
}

func (MS *MemStorage) Close() error {
	return nil
}

func (MS *MemStorage) AddUser(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	MS.mu.Lock()
	defer MS.mu.Unlock()
	for _, v := range MS.userCash {
		if v.Login == user.Login {
			return constans.ErrorNoUNIQUE
		}
	}
	MS.userCash[uuid.New()] = user
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (MS *MemStorage) AuthenticationUser(_ context.Context, user models.User) (bool, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	for _, v := range MS.userCash {
		if v.Login == user.Login && v.Password == user.Password {
			return true, nil
		}
	}
	return false, nil
}

func (MS *MemStorage) GetOrder(_ context.Context, numberOrder string) (models.Order, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()

	return MS.orderCash[numberOrder], nil
}

func (MS *MemStorage) GetManyOrders(_ context.Context, userLogin string) ([]models.Order, error) {
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

func (MS *MemStorage) AddOrder(_ context.Context, numberOrder string, order models.Order) error {
	MS.mu.Lock()
	defer MS.mu.Unlock()
	if MS.orderCash[numberOrder].NumberOrder == numberOrder {
		return constans.ErrorNoUNIQUE
	}
	MS.orderCash[numberOrder] = order
	return nil
}

func (MS *MemStorage) GetOrdersByProcess(_ context.Context) ([]models.Order, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	var orders []models.Order
	for _, v := range MS.orderCash {
		if v.Status == constans.OrderStatusPROCESSING ||
			v.Status == constans.OrderStatusNEW ||
			v.Status == constans.OrderStatusREGISTERED {
			orders = append(orders, v)
		}
	}
	return orders, nil
}

func (MS *MemStorage) UpdateOrder(_ context.Context, loyaltyPoint models.LoyaltyPoint) error {
	MS.mu.Lock()
	defer MS.mu.Unlock()
	order := MS.orderCash[loyaltyPoint.NumberOrder]
	order.Status, order.Accrual = loyaltyPoint.Status, loyaltyPoint.Accrual
	MS.orderCash[loyaltyPoint.NumberOrder] = order
	return nil
}

func (MS *MemStorage) GetUserBalance(_ context.Context, userLogin string) (float64, float64, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	pointsSUM := 0.0
	pointsSPEND := 0.0
	for _, order := range MS.orderCash {
		if order.UserLogin == userLogin {
			pointsSUM += order.Accrual
		}
	}
	for _, withdraw := range MS.withdrawCash {
		if withdraw.UserLogin == userLogin {
			pointsSPEND += withdraw.Sum
		}
	}
	return pointsSUM, pointsSPEND, nil
}

func (MS *MemStorage) AddWithdraw(ctx context.Context, withdraw models.Withdraw) error {
	orderSum, wSum, err := MS.GetUserBalance(ctx, withdraw.UserLogin)
	if err != nil {
		return err
	}
	if orderSum < wSum+withdraw.Sum {
		return constans.ErrorStatusShortfallAccount
	}
	MS.mu.Lock()
	defer MS.mu.Unlock()
	MS.withdrawCash[uuid.New()] = withdraw
	return nil
}

func (MS *MemStorage) GetManyWithdraws(_ context.Context, userLogin string) ([]models.Withdraw, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	var withdraws []models.Withdraw
	for _, v := range MS.withdrawCash {
		if v.UserLogin == userLogin {
			withdraws = append(withdraws, v)
		}
	}
	return withdraws, nil
}
