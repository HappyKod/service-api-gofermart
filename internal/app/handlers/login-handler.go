package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"HappyKod/service-api-gofermart/internal/utils"
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoginHandler Аутентификация пользователя
// Хендлер: POST /api/user/login.
//
// Аутентификация производится по паре логин/пароль.
//
// Формат запроса:
//
// POST /api/user/login HTTP/1.1
// Content-Type: application/json
// ...
//
//	{
//		"login": "<login>",
//		"password": "<password>"
//	}
//
// Возможные коды ответа:
//
// 200 — пользователь успешно аутентифицирован;
// 400 — неверный формат запроса;
// 401 — неверная пара логин/пароль;
// 500 — внутренняя ошибка сервера.
func LoginHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constans.TimeOutRequest)
	defer cancel()
	if !utils.ValidContentType(c, "application/json") {
		return
	}
	log := container.GetLog()
	storage := container.GetStorage()
	var user models.User
	if err := c.Bind(&user); err != nil {
		log.Error(constans.ErrorUnmarshalBody, zap.Error(err))
		c.String(http.StatusInternalServerError, constans.ErrorUnmarshalBody)
		return
	}
	log.Debug("авторизация пользователя", zap.Any("user", user))
	if user.Login == "" || user.Password == "" {
		log.Debug("не валидные логин или пароль", zap.Any("user", user))
		c.String(http.StatusBadRequest, "не валидные логин или пароль")
		return
	}
	authenticationUser, err := storage.AuthenticationUser(ctx, user)
	if err != nil {
		log.Error(constans.ErrorWorkDataBase, zap.Error(err))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	if !authenticationUser {
		log.Debug("пароль или логин не верный", zap.Any("user", user))
		c.String(http.StatusUnauthorized, "пароль или логин не верный")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 100)),
			IssuedAt:  jwt.At(time.Now())},
		Login: user.Login,
	})
	log.Debug("пользователь успешно авторизовался",
		zap.Any("user", user),
		zap.Any("token", token))
	accessToken, err := token.SignedString([]byte(container.GetConfig().SecretKey))
	if err != nil {
		log.Error("ошибка генерация токена", zap.Error(err))
		c.String(http.StatusInternalServerError, "ошибка генерация токена")
		return
	}
	c.Header("Authorization", "Bearer "+accessToken)
}
