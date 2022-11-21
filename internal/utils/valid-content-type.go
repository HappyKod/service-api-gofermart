package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidContentType Валидация ContentType
func ValidContentType(c *gin.Context, ContentType string) bool {
	if c.GetHeader("Content-Type") != ContentType {
		c.String(http.StatusBadRequest, "не верный заголовок ожидалось %s", ContentType)
		return false
	}
	return true
}
