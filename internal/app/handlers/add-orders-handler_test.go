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

func TestAddUserOrders(t *testing.T) {
	type want struct {
		responseCode int
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
			name:          "регистрация номера заказа без авторизации",
			requestPath:   "/api/user/orders",
			requestMethod: http.MethodPost,
			requestBody:   "",
			requestHeader: [2]string{"Content-Type", "text/plain"},
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
			name:          "регистрация повторного номера заказа c авторизацией",
			requestPath:   "/api/user/orders",
			requestMethod: http.MethodPost,
			requestBody:   "4561261212345467",
			requestHeader: [2]string{"Content-Type", "text/plain"},
			want: want{
				responseCode: http.StatusOK,
			},
		},
		{
			name:          "регистрация номера заказа c авторизацией не правильный Content-Type",
			requestPath:   "/api/user/orders",
			requestMethod: http.MethodPost,
			requestBody:   "4561261212345467",
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name:          "регистрация номера заказа c авторизацией не верный код заказа",
			requestPath:   "/api/user/orders",
			requestMethod: http.MethodPost,
			requestBody:   "4561261212345464",
			requestHeader: [2]string{"Content-Type", "text/plain"},
			want: want{
				responseCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name:          "регистрация номера заказа c авторизацией",
			requestPath:   "/api/user/orders",
			requestMethod: http.MethodPost,
			requestBody:   "4111111111111111",
			requestHeader: [2]string{"Content-Type", "text/plain"},
			want: want{
				responseCode: http.StatusAccepted,
			},
		},
		{
			name:          "регистрация номера заказа c авторизацией",
			requestPath:   "/api/user/orders",
			requestMethod: http.MethodPost,
			requestBody:   "79927398713",
			requestHeader: [2]string{"Content-Type", "text/plain"},
			want: want{
				responseCode: http.StatusAccepted,
			},
		},
		{
			name:          "регистрация пользователя",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin2",
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
				"login":"admin2",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusOK,
			},
		},
		{
			name:          "регистрация повторного номера заказа c авторизацией",
			requestPath:   "/api/user/orders",
			requestMethod: http.MethodPost,
			requestBody:   "79927398713",
			requestHeader: [2]string{"Content-Type", "text/plain"},
			want: want{
				responseCode: http.StatusConflict,
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
			assert.Equal(t, tt.want.responseCode, w.Code)
		})
	}
}
