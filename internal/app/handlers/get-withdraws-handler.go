package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetUserWithdraws Получение информации о выводе средств
// Хендлер: GET /api/user/balance/withdrawals.
//
// Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым. Формат даты — RFC3339.
//
// Формат запроса:
//
// GET /api/user/withdrawals HTTP/1.1
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
// [
//
//	{
//	    "order": "2377225624",
//	    "sum": 500,
//	    "processed_at": "2020-12-09T16:09:57+03:00"
//	}
//
// ]
// 204 - нет ни одного списания.
//
// 401 — пользователь не авторизован.
//
// 500 — внутренняя ошибка сервера.
func GetUserWithdraws(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constans.TimeOutRequest)
	defer cancel()
	log := container.GetLog()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	log.Debug("поступил запрос на показ списаний",
		zap.String("loginUser", user))
	orders, err := storage.GetManyWithdraws(ctx, user)
	if err != nil {
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "GetManyWithdraws"))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	if len(orders) == 0 {
		log.Debug("нет данных для ответа", zap.String("loginUser", user))
		c.String(http.StatusNoContent, "нет данных для ответа")
		return
	}
	c.JSON(http.StatusOK, orders)
}
