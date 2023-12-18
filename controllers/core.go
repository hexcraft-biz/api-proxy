package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/drawbridge/config"
)

type Controller struct {
	*config.Config
}

func NewController(cfg *config.Config) *Controller {
	return &Controller{
		Config: cfg,
	}
}

type Common struct {
	*Controller
}

func NewCommon(cfg *config.Config) *Common {
	return &Common{
		Controller: NewController(cfg),
	}
}

func (ctrl *Common) NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
	}
}

func (ctrl *Common) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}
}
