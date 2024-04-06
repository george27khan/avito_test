package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Привет, Stepik!",
	})
}
