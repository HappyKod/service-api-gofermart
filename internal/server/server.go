package server

import (
	"HappyKod/service-api-gofermart/internal/models"
	"github.com/gin-gonic/gin"
	"log"
)

// NewServer создания сервера с настройками
func NewServer(r *gin.Engine, cfg models.Config) {
	log.Fatalln(r.Run(cfg.Address))
}
