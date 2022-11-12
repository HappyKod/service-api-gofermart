package main

import (
	"HappyKod/service-api-gofermart/internal/container"
	"HappyKod/service-api-gofermart/internal/handlers"
	"HappyKod/service-api-gofermart/internal/models"
	"HappyKod/service-api-gofermart/internal/server"
	"flag"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"log"
)

var cfg models.Config

func init() {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "адрес запуска HTTP-сервера")
	flag.StringVar(&cfg.DataBaseURI, "d", cfg.DataBaseURI, "строка с адресом подключения к БД")
	flag.StringVar(&cfg.AccrualAddress, "r", cfg.AccrualAddress, "адрес системы расчёта начислений")
}
func main() {
	var zapLogger *zap.Logger
	var err error
	if err = env.Parse(&cfg); err != nil {
		log.Fatalln("ошибка считывания конфига", zap.Error(err))
	}
	flag.Parse()
	if cfg.ReleaseMOD {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalln(err)
	}
	zapLogger.Info("считана следующая конфигурация",
		zap.String("AddressServer", cfg.Address),
		zap.String("AccrualAddress", cfg.AccrualAddress),
		zap.Bool("ReleaseMOD", cfg.ReleaseMOD),
	)
	zapLogger.Debug("полная конфигурация", zap.Any("config", cfg))
	if err = container.BuildContainer(cfg, zapLogger); err != nil {
		zapLogger.Fatal("ошибка запуска Di контейнера", zap.Error(err))
	}
	r := handlers.Router(cfg)
	server.NewServer(r, cfg)
}
