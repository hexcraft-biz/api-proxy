package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
	}
}

func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}
}
