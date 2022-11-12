package storage

type GoferStorage interface {
	Ping() error
	Close() error
}
