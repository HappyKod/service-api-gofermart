package handlers

import (
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
	gUser := r.Group("/api/user")
	gUserBalance := gUser.Group("/balance")
	{
		gUser.POST("/register", func(context *gin.Context) {

		})
		gUser.POST("/login", func(context *gin.Context) {})
		gUser.POST("/orders", func(context *gin.Context) {})
		gUser.GET("/orders", func(context *gin.Context) {})
		gUserBalance.GET("/", func(context *gin.Context) {})
		gUserBalance.POST("/withdraw", func(context *gin.Context) {})
		gUserBalance.GET("/withdrawals", func(context *gin.Context) {})
	}
	return r
}
