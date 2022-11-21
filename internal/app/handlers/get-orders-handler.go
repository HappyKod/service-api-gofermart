package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetUserOrders Получение списка загруженных номеров заказов
// Хендлер: GET /api/user/orders.
//
// Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым. Формат даты — RFC3339.
//
// Доступные статусы обработки расчётов:
//
// NEW — заказ загружен в систему, но не попал в обработку;
// PROCESSING — вознаграждение за заказ рассчитывается;
// INVALID — система расчёта вознаграждений отказала в расчёте;
// PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
// Формат запроса:
//
// GET /api/user/orders HTTP/1.1
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
//		{
//	       "number": "9278923470",
//	       "status": "PROCESSED",
//	       "accrual": 500,
//	       "uploaded_at": "2020-12-10T15:15:45+03:00"
//	   },
//	   {
//	       "number": "12345678903",
//	       "status": "PROCESSING",
//	       "uploaded_at": "2020-12-10T15:12:01+03:00"
//	   },
//	   {
//	       "number": "346436439",
//	       "status": "INVALID",
//	       "uploaded_at": "2020-12-09T16:09:53+03:00"
//	   }
//
// ]
// 204 — нет данных для ответа.
//
// 401 — пользователь не авторизован.
//
// 500 — внутренняя ошибка сервера.
func GetUserOrders(c *gin.Context) {
	log := container.GetLog()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	log.Debug("поступил запрос на показ заказов",
		zap.String("loginUser", user))

	orders, err := storage.GetManyOrders(user)
	if err != nil {
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "GetManyOrders"))
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
