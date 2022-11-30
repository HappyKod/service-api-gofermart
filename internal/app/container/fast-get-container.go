package container

import (
	"HappyKod/service-api-gofermart/internal/models"
	"HappyKod/service-api-gofermart/internal/storage"

	"go.uber.org/zap"
)

// GetLog
func GetLog() *zap.Logger {
	return DiContainer.Get("zap-logger").(*zap.Logger)
}

func GetStorage() storage.GoferStorage {
	return DiContainer.Get("storage").(storage.GoferStorage)
}

func GetConfig() models.Config {
	return DiContainer.Get("server-config").(models.Config)
}
