package service

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// CalculationLoyaltyPoints Сервис по обновлению
// показателей статуса заказа
func CalculationLoyaltyPoints(ctx context.Context) error {
	storage := container.GetStorage()
	cfg := container.GetConfig()
	log := container.GetLog()
	ctxStorage, cancelStorage := context.WithTimeout(ctx, constans.TimeOutStorage)
	orders, err := storage.GetOrdersByProcess(ctxStorage)
	cancelStorage()
	if err != nil {
		return err
	}
	for _, order := range orders {
		ctx, cancel := context.WithTimeout(ctx, constans.TimeOutRequest)
		joinPath, errJoin := url.JoinPath(cfg.AccrualAddress, "/api/orders/", order.NumberOrder)
		if errJoin != nil {
			cancel()
			return errJoin
		}
		r, errGet := http.Get(joinPath)
		if errGet != nil {
			cancel()
			return errGet
		}
		switch r.StatusCode {
		case http.StatusTooManyRequests:
			log.Error("TooManyRequests")
			time.Sleep(constans.TimeSleepTooManyRequests)
		case http.StatusInternalServerError:
			body, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				cancel()
				return errRead
			}
			log.Error("StatusInternalServerError", zap.String("body", string(body)))
			errClose := r.Body.Close()
			if errClose != nil {
				cancel()
				return errClose
			}
		case http.StatusOK:
			body, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				cancel()
				return errRead
			}
			log.Debug("body", zap.String("body", string(body)))
			var loyaltyPoint models.LoyaltyPoint
			if errUnmarshal := json.Unmarshal(body, &loyaltyPoint); errUnmarshal != nil {
				cancel()
				return errUnmarshal
			}
			errClose := r.Body.Close()
			if errClose != nil {
				cancel()
				return errClose
			}
			if loyaltyPoint.Status == constans.OrderStatusREGISTERED {
				loyaltyPoint.Status = constans.OrderStatusPROCESSING
			}
			if order.Status != loyaltyPoint.Status {
				if err = storage.UpdateOrder(ctx, loyaltyPoint); err != nil {
					cancel()
					return err
				}
			}
		}
		cancel()
	}
	return nil
}
