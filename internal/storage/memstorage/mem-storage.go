package memstorage

import (
	"HappyKod/service-api-gofermart/internal/models"
	"github.com/google/uuid"
	"sync"
)

type MemStorage struct {
	userCash map[uuid.UUID]models.User
	mu       *sync.RWMutex
}

func New() (*MemStorage, error) {
	return &MemStorage{
		userCash: make(map[uuid.UUID]models.User),
		mu:       new(sync.RWMutex),
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
