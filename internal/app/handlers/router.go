package handlers

import (
	"HappyKod/service-api-gofermart/internal/app/middleware"
	"HappyKod/service-api-gofermart/internal/models"
	"github.com/gin-gonic/gin"
)

// Router указание маршрутов хендрлеров
func Router(cfg models.Config) *gin.Engine {
	if cfg.ReleaseMOD {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.JwtValid())
	gUser := r.Group("/api/user")
	gUserBalance := gUser.Group("/balance")

	{
		gUser.POST("/register", RegisterHandler)
		gUser.POST("/login", LoginHandler)
		gUser.POST("/orders", AddUserOrders)
		gUser.GET("/orders", GetUserOrders)
		gUserBalance.GET("/", func(context *gin.Context) {})
		gUserBalance.POST("/withdraw", func(context *gin.Context) {})
		gUserBalance.GET("/withdrawals", func(context *gin.Context) {})
	}
	return r
}
