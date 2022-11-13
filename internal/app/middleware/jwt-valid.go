package middleware

import (
	"HappyKod/service-api-gofermart/internal/app/container"
	"HappyKod/service-api-gofermart/internal/models"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/auth/pkg/auth"
	"net/http"
	"strings"
)

func JwtValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/user/register" || c.Request.URL.Path == "/api/user/login" {
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if headerParts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_, err := parseToken(headerParts[1],
			[]byte(container.DiContainer.Get("server-config").(models.Config).SecretKey),
		)
		if err != nil {
			status := http.StatusBadRequest
			if err == auth.ErrInvalidAccessToken {
				status = http.StatusUnauthorized
			}

			c.AbortWithStatus(status)
			return
		}
		fmt.Println(c.Request.URL.Path)
	}
}

func parseToken(accessToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
		return claims.Username, nil
	}

	return "", auth.ErrInvalidAccessToken
}
