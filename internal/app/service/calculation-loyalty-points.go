package service

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"time"
)

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
			var loyaltyPoints []models.LoyaltyPoint
			if errUnmarshal := json.Unmarshal(body, &loyaltyPoints); errUnmarshal != nil {
				return errUnmarshal
			}
			errClose := r.Body.Close()
			if errClose != nil {
				return errClose
			}
			for _, loyaltyPoint := range loyaltyPoints {
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
	}
	return nil
}
