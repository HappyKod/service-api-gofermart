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
	"strconv"
	"time"
)

// AddWithdraw Запрос на списание средств
// Хендлер: POST /api/user/balance/withdraw
//
// Хендлер доступен только авторизованному пользователю.
// Номер заказа представляет собой гипотетический номер нового заказа пользователя в счет оплаты которого списываются баллы.
//
// Примечание: для успешного списания достаточно успешной регистрации запроса,
// никаких внешних систем начисления не предусмотрено и не требуется реализовывать.
//
// Формат запроса:
//
// POST /api/user/balance/withdraw HTTP/1.1
// Content-Type: application/json
//
//	{
//		"order": "2377225624",
//	   "sum": 751
//	}
//
// Здесь order — номер заказа, а sum — сумма баллов к списанию в счёт оплаты.
//
// Возможные коды ответа:
//
// 200 — успешная обработка запроса;
// 401 — пользователь не авторизован;
// 402 — на счету недостаточно средств;
// 422 — неверный номер заказа;
// 500 — внутренняя ошибка сервера.
func AddWithdraw(c *gin.Context) {
	if !utils.ValidContentType(c, "application/json") {
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
	var withdraw models.Withdraw
	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		log.Error(constans.ErrorUnmarshalBody, zap.Error(err))
		c.String(http.StatusInternalServerError, constans.ErrorUnmarshalBody)
		return
	}
	withdraw.ProcessedAT, withdraw.UserLogin = time.Now(), user
	log.Debug("поступил запрос на списание средств",
		zap.Any("withdraw", withdraw),
		zap.String("loginUser", user))
	numberOrder, err := strconv.Atoi(withdraw.NumberOrder)
	if err != nil {
		log.Debug("ошибка преобразования номера заказа", zap.Any("withdraw", withdraw))
		c.String(http.StatusInternalServerError, "ошибка преобразования номера заказа")
		return
	}
	if !utils.ValidLuhn(numberOrder) {
		log.Debug(constans.ErrorNumberValidLuhn, zap.Error(err), zap.Int("numberOrder", numberOrder))
		c.String(http.StatusUnprocessableEntity, constans.ErrorNumberValidLuhn)
		return
	}

	sum, spent, err := storage.GetUserBalance(user)
	if err != nil {
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "GetUserBalance"))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	current := sum - spent
	if current-withdraw.Sum < 0 {
		log.Debug("на счету недостаточно средств", zap.String("user", user),
			zap.Float64("current", current),
			zap.Float64("current", withdraw.Sum),
		)
		c.String(http.StatusPaymentRequired, "на счету недостаточно средств")
		return
	}
	err = storage.AddWithdraw(withdraw)
	if err != nil {
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "AddWithdraw"))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	log.Debug("списание совершено", zap.Any("withdraw", withdraw))
	c.String(http.StatusOK, "списание совершено")
}
