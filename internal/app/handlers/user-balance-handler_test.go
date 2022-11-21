package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/models"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.uber.org/zap"
)

func TestUserBalance(t *testing.T) {
	type want struct {
		responseCode int
		body         string
	}
	tests := []struct {
		name          string
		requestPath   string
		requestMethod string
		requestBody   string
		requestHeader [2]string
		want          want
		cfg           models.Config
	}{
		{
			name:          "просмотр баланса без авторизации",
			requestPath:   "/api/user/balance",
			requestMethod: http.MethodGet,
			requestBody:   "",
			requestHeader: [2]string{"Content-Length", "0"},
			want: want{
				responseCode: http.StatusUnauthorized,
			},
		},
		{
			name:          "регистрация пользователя",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin1",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusPermanentRedirect,
			},
		},
		{
			name:          "авторизация пользователя",
			requestPath:   "/api/user/login",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin1",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusOK,
			},
		},
		{
			name:          "просмотр баланса с авторизацией без заказов",
			requestPath:   "/api/user/balance",
			requestMethod: http.MethodGet,
			requestBody:   "",
			requestHeader: [2]string{"Content-Length", "0"},
			want: want{
				responseCode: http.StatusOK,
				body:         `{"current":0,"withdrawn":0}`,
			},
		},
		{
			name:          "регистрация номера заказа c авторизацией",
			requestPath:   "/api/user/orders",
			requestMethod: http.MethodPost,
			requestBody:   "4561261212345467",
			requestHeader: [2]string{"Content-Type", "text/plain"},
			want: want{
				responseCode: http.StatusAccepted,
			},
		},
		{
			name:          "просмотр баланса с авторизацией без заказов",
			requestPath:   "/api/user/balance",
			requestMethod: http.MethodGet,
			requestBody:   "",
			requestHeader: [2]string{"Content-Length", "0"},
			want: want{
				responseCode: http.StatusOK,
				body:         `{"current":100,"withdrawn":0}`,
			},
		},
		{
			name:          "регистрация списание бонусов c авторизацией",
			requestPath:   "/api/user/balance/withdraw",
			requestMethod: http.MethodPost,
			requestBody:   `{"order": "4111111111111111","sum": 10}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusOK,
			},
		},
		{
			name:          "просмотр баланса с авторизацией без заказов",
			requestPath:   "/api/user/balance",
			requestMethod: http.MethodGet,
			requestBody:   "",
			requestHeader: [2]string{"Content-Length", "0"},
			want: want{
				responseCode: http.StatusOK,
				body:         `{"current":90,"withdrawn":10}`,
			},
		},
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	err = container.BuildContainer(models.Config{}, logger)
	if err != nil {
		t.Fatal(err)
	}
	var bearer string
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := Router(models.Config{})
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Add(tt.requestHeader[0], tt.requestHeader[1])
			req.Header.Add("Authorization", bearer)
			router.ServeHTTP(w, req)
			if tt.requestPath == "/api/user/login" {
				bearer = w.Header().Get("Authorization")
			}
			if tt.requestPath == "/api/user/orders" {
				err = container.GetStorage().UpdateOrder(models.LoyaltyPoint{
					Status:      "PROCESSED",
					NumberOrder: tt.requestBody,
					Accrual:     100.0})
				if err != nil {
					t.Fatal(err)
				}
			}
			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, w.Body.String())
			}
			assert.Equal(t, tt.want.responseCode, w.Code)
		})
	}
}
