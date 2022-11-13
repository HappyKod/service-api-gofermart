package storage

import "HappyKod/service-api-gofermart/internal/models"

type GoferStorage interface {
	Ping() error
	Close() error
	UniqLoginUser(login string) (bool, error)
	AddUser(user models.User) error
	AuthenticationUser(user models.User) (bool, error)
}
