package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/models"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.uber.org/zap"
)

func TestLoginHandler(t *testing.T) {
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
			name:          "авторизация пользователя без регистрации",
			requestPath:   "/api/user/login",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin1",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
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
			name:          "авторизация пользователя с неверный паролем",
			requestPath:   "/api/user/login",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin1",
				"password": "admin1"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusUnauthorized,
			},
		},
		{
			name:          "авторизация пользователя с неверный логином",
			requestPath:   "/api/user/login",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin",
				"password": "admin"
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusUnauthorized,
			},
		},
		{
			name:          "авторизация пользователя без правильных заголовков",
			requestPath:   "/api/user/login",
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
			name:          "авторизация пользователя с пустыми данными",
			requestPath:   "/api/user/login",
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
			name:          "авторизация пользователя с пустыми данными",
			requestPath:   "/api/user/login",
			requestMethod: http.MethodPost,
			requestBody: `{
				"login":"admin1",
				"password": ""
			}`,
			requestHeader: [2]string{"Content-Type", "application/json"},
			want: want{
				responseCode: http.StatusBadRequest,
			},
		},
		{
			name:          "авторизация пользователя с пустыми данными",
			requestPath:   "/api/user/login",
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
			log.Println(w.Header(), w.Code)
			assert.Equal(t, tt.want.responseCode, w.Code)
		})
	}
}
