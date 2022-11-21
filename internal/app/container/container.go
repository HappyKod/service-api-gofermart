package container

import (
	"HappyKod/service-api-gofermart/internal/models"
	"HappyKod/service-api-gofermart/internal/storage/memstorage"

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
	storage, err := memstorage.New()
	if err != nil {
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
		Build: func(ctn di.Container) (interface{}, error) { return storage, nil }}); err != nil {
		return err
	}
	DiContainer = builder.Build()
	return nil
}
