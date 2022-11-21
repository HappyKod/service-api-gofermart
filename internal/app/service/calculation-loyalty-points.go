package service

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// CalculationLoyaltyPoints Сервис по обновлению
// показателей статуса заказа
func CalculationLoyaltyPoints() error {
	storage := container.GetStorage()
	cfg := container.GetConfig()
	log := container.GetLog()
	orders, err := storage.GetOrdersByProcess()
	if err != nil {
		return err
	}
	for _, order := range orders {
		joinPath, errJoin := url.JoinPath(cfg.AccrualAddress, "/api/orders/", order.NumberOrder)
		if errJoin != nil {
			return errJoin
		}
		r, errGet := http.Get(joinPath)
		if errGet != nil {
			return errGet
		}
		switch r.StatusCode {
		case http.StatusTooManyRequests:
			log.Error("TooManyRequests")
			time.Sleep(constans.TimeSleepTooManyRequests * time.Second)
		case http.StatusInternalServerError:
			body, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				return errRead
			}
			log.Error("StatusInternalServerError", zap.String("body", string(body)))
			errClose := r.Body.Close()
			if errClose != nil {
				return errClose
			}
		case http.StatusOK:
			body, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				return errRead
			}
			log.Debug("body", zap.String("body", string(body)))
			var loyaltyPoint models.LoyaltyPoint
			if errUnmarshal := json.Unmarshal(body, &loyaltyPoint); errUnmarshal != nil {

				return errUnmarshal
			}
			errClose := r.Body.Close()
			if errClose != nil {
				return errClose
			}
			if loyaltyPoint.Status == "REGISTERED" {
				loyaltyPoint.Status = "PROCESSING"
			}
			if order.Status != loyaltyPoint.Status {
				if err = storage.UpdateOrder(loyaltyPoint); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
