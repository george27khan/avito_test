package handlers

import (
	"avito_test/pkg/auth"
	usr "avito_test/pkg/postgres_db/user"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GetToken функция получения токена
func GetToken(c *gin.Context) {
	ctx := context.Background()
	userName, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	fmt.Println("userName, password, ok", userName, password, ok)
	user, err := usr.Get(ctx, userName)
	fmt.Println("user, err", user, err)

	if err != nil || !user.VerifyPassword(password) {
		c.Status(http.StatusUnauthorized)
		return
	} else {
		if token, err := auth.GetToken(user.UserName, user.Password, user.IsAdmin); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		} else {
			c.Writer.Header().Set("token", token)
			c.Status(http.StatusOK)
			return
		}
	}
}

// Auth функция аутентификации пользователя
func Auth(c *gin.Context) {
	token := c.GetHeader("token")
	fmt.Println("token", token)
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := auth.ParseToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println(claims)
	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("is_admin", claims["is_admin"])
}

// AuthMiddleware middleware аутентификации
func AuthMiddleware() gin.HandlerFunc {
	return Auth
}
