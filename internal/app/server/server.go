package server

import (
	"HappyKod/service-api-gofermart/internal/models"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// NewServer создания сервера с настройками
func NewServer(r *gin.Engine, cfg models.Config) {
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ошибка: %s\n\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("экстренное выключение сервиса", err)
	}
	log.Println("сервис выключен")
}
