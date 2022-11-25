package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/constans"
	"HappyKod/service-api-gofermart/internal/models"
	"HappyKod/service-api-gofermart/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), constans.TimeOutRequest)
	defer cancel()
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
	err = storage.AddWithdraw(ctx, withdraw)
	if err != nil {
		if errors.Is(err, constans.StatusShortfallAccount) {
			c.String(http.StatusPaymentRequired, constans.StatusShortfallAccount.Error())
			return
		}
		log.Error(constans.ErrorWorkDataBase, zap.Error(err), zap.String("func", "AddWithdraw"))
		c.String(http.StatusInternalServerError, constans.ErrorWorkDataBase)
		return
	}
	log.Debug("списание совершено", zap.Any("withdraw", withdraw))
	c.String(http.StatusOK, "списание совершено")
}
