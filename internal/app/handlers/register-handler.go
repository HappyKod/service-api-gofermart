package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"HappyKod/service-api-gofermart/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterHandler Регистрация пользователя
// Хендлер: POST /api/user/register.
//
// Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным. После успешной регистрации должна происходить автоматическая аутентификация пользователя.
//
// Формат запроса:
//
// POST /api/user/register HTTP/1.1
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
// 200 — пользователь успешно зарегистрирован и аутентифицирован;
// 400 — неверный формат запроса;
// 409 — логин уже занят;
// 500 — внутренняя ошибка сервера.
func RegisterHandler(c *gin.Context) {
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
	log.Debug("регистрация пользователя", zap.Any("user", user))
	if user.Login == "" || user.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	uniq, err := storage.UniqLoginUser(user.Login)
	if err != nil {
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "UniqLoginUser"))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	if !uniq {
		log.Debug("пользователь с таким логином уже есть", zap.Any("user", user))
		c.String(http.StatusConflict, "пользователь с таким логином уже есть")
		return
	}
	err = storage.AddUser(user)
	if err != nil {
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "AddUser"))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	log.Debug("пользователь успешно зарегистрирован", zap.Any("user", user))
	c.Redirect(http.StatusPermanentRedirect, "/api/user/login")
}
