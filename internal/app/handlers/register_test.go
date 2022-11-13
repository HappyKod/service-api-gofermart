package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/models"
	"bytes"
	"github.com/go-playground/assert/v2"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
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
			name:          "повторная попытка регистрация пользователя",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin1",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusConflict,
			},
		},
		{
			name:          "регистрация пользователя с пустыми данными",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name:          "регистрация пользователя с пустыми данными",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin123",
				"password": ""
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name:          "регистрация пользователя с пустыми данными",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"",
				"password": ""
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name:          "регистрация пользователя без правильных заголовков",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin1",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "text/plain"},
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name:          "регистрация пользователя добавочного",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin10",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusPermanentRedirect,
			},
		},
		{
			name:          "регистрация пользователя добавочного",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin11",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusPermanentRedirect,
			},
		},
		{
			name:          "регистрация пользователя добавочного",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin12",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusPermanentRedirect,
			},
		},
		{
			name:          "регистрация пользователя добавочного",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin13",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusPermanentRedirect,
			},
		},
		{
			name:          "регистрация пользователя добавочного",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin14",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusPermanentRedirect,
			},
		},
		{
			name:          "регистрация пользователя добавочного",
			requestPath:   "/api/user/register",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin15",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusPermanentRedirect,
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := Router(models.Config{})
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Add(tt.requestHeader[0], tt.requestHeader[1])
			router.ServeHTTP(w, req)
			log.Println(w.Body.String(), w.Code)
			assert.Equal(t, tt.want.responseCode, w.Code)
		})
	}
}
