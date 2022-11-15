package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			lenContent := c.GetHeader("Content-Length")
			if lenContent != "0" {
				c.String(http.StatusBadRequest, "не соответствие заголовка Content-Length:0")
				return
			}
		}
		c.Next()
	}
}
