package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserBalance Получение текущего баланса пользователя
// Хендлер: GET /api/user/balance.
//
// Хендлер доступен только авторизованному пользователю.
// В ответе должны содержаться данные о текущей сумме баллов лояльности,
// а также сумме использованных за весь период регистрации баллов.
//
// Формат запроса:
//
// GET /api/user/balance HTTP/1.1
// Content-Length: 0
// Возможные коды ответа:
//
// 200 — успешная обработка запроса.
//
// Формат ответа:
//
// 200 OK HTTP/1.1
// Content-Type: application/json
// ...
//
//	{
//		"current": 500.5,
//		"withdrawn": 42
//	}
//
// 401 — пользователь не авторизован.
//
// 500 — внутренняя ошибка сервера.
func UserBalance(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constans.TimeOutRequest)
	defer cancel()
	log := container.GetLog()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	log.Debug("поступил запрос на проверку баланса",
		zap.String("loginUser", user))

	sum, spent, err := storage.GetUserBalance(ctx, user)
	if err != nil {
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "GetUserBalance"))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	response := map[string]float64{
		"current":   sum - spent,
		"withdrawn": spent,
	}
	log.Debug("баланс пользователя", zap.String("loginUser", user),
		zap.Float64("sum", sum),
		zap.Float64("spent", spent),
		zap.Float64("current", sum-spent),
	)
	c.JSONP(http.StatusOK, response)
}
