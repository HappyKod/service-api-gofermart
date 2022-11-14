package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"HappyKod/service-api-gofermart/internal/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

// AddUserOrders Загрузка номера заказа
// Хендлер: POST /api/user/orders.
//
// Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
//
// Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.
//
// Формат запроса:
//
// POST /api/user/orders HTTP/1.1
// Content-Type: text/plain
// ...
//
// 12345678903
// Возможные коды ответа:
//
// 200 — номер заказа уже был загружен этим пользователем;
// 202 — новый номер заказа принят в обработку;
// 400 — неверный формат запроса;
// 401 — пользователь не аутентифицирован;
// 409 — номер заказа уже был загружен другим пользователем;
// 422 — неверный формат номера заказа;
// 500 — внутренняя ошибка сервера.
func AddUserOrders(c *gin.Context) {
	if !utils.ValidContentType(c, "text/plain") {
		return
	}
	log := container.GetLog()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(constans.ErrorReadBody, zap.Error(err))
		c.String(http.StatusInternalServerError, constans.ErrorReadBody)
		return
	}
	var numberOrder int
	err = json.Unmarshal(body, &numberOrder)
	if err != nil {
		log.Error(constans.ErrorUnmarshalBody, zap.Error(err))
		c.String(http.StatusInternalServerError, constans.ErrorUnmarshalBody)
		return
	}
	log.Debug("поступил номер заказа",
		zap.Int("numberOrder", numberOrder),
		zap.String("loginUser", user))

	if !utils.ValidLuhn(numberOrder) {
		log.Debug(constans.ErrorNumberValidLuhn, zap.Error(err), zap.Int("numberOrder", numberOrder))
		c.String(http.StatusUnprocessableEntity, constans.ErrorNumberValidLuhn)
		return
	}
	err = storage.AddOrder(numberOrder,
		models.Order{
			NumberOrder: numberOrder,
			UserLogin:   user,
			Status:      constans.OrderStatusPROCESSING,
			Created:     time.Now(),
		})
	if err != nil {
		if err.Error() == constans.ErrorNoUNIQUE {
			order, errGet := storage.GetOrder(numberOrder)
			if errGet != nil {
				log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "GetOrder"))
				c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
				return
			}
			if order.UserLogin == user {
				log.Debug("номер заказа уже был загружен этим пользователем", zap.Any("order", order))
				c.String(http.StatusOK, "номер заказа уже был загружен этим пользователем")
				return
			}
			log.Debug("номер заказа уже был загружен другим пользователем", zap.Any("order", order))
			c.String(http.StatusConflict, "номер заказа уже был загружен другим пользователем")
			return
		}
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "AddOrder"))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	log.Debug("новый номер заказа принят в обработку", zap.Any("number_order", numberOrder))
	c.String(http.StatusAccepted, "новый номер заказа принят в обработку")
}
