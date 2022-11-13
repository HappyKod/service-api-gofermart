package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidContentType(c *gin.Context, ContentType string) bool {
	if c.GetHeader("Content-Type") != ContentType {
		c.String(http.StatusBadRequest, "не верный заголовок ожидалось %s", ContentType)
		return false
	}
	return true
}
