package main

import (
	"avito_test/cmd/web/banner_handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/user_banner", banner_handler.Handler)
	r.GET("/banner", banner_handler.Handler)
	r.POST("/banner", banner_handler.Handler)
	r.PATCH("/banner/:id", banner_handler.Handler)
	r.DELETE("/banner/:id", banner_handler.Handler)

	r.Run("localhost:8080") //63342
}
