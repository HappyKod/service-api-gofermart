package container

import (
	"HappyKod/service-api-gofermart/internal/models"
	"HappyKod/service-api-gofermart/internal/storage"
	"HappyKod/service-api-gofermart/internal/storage/memstorage"
	"HappyKod/service-api-gofermart/internal/storage/pgstorage"

	"github.com/sarulabs/di"
	"go.uber.org/zap"
)

var DiContainer di.Container

// BuildContainer Создание контейнера
// с подключениями, и необходимыми сущностями по всему коду
// и присваивание DiContainer
func BuildContainer(cfg models.Config, logger *zap.Logger) error {
	builder, err := di.NewBuilder()
	if err != nil {
		return err
	}
	var goferStorage storage.GoferStorage
	if cfg.DataBaseURI != "" {
		goferStorage, err = pgstorage.New(cfg.DataBaseURI)
		if err != nil {
			return err
		}
	} else {
		goferStorage, err = memstorage.New()
		if err != nil {
			return err
		}
	}
	if err = goferStorage.Ping(); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "server-config",
		Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }}); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "zap-logger",
		Build: func(ctn di.Container) (interface{}, error) { return logger, nil }}); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "storage",
		Build: func(ctn di.Container) (interface{}, error) { return goferStorage, nil }}); err != nil {
		return err
	}
	DiContainer = builder.Build()
	return nil
}
