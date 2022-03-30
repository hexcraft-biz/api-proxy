package route

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/karmaksana-io/api-proxy/common"
	"github.com/karmaksana-io/api-proxy/config"
	"github.com/karmaksana-io/api-proxy/controller"
	"github.com/karmaksana-io/api-proxy/middleware"
)

// all the routes are defined here
func NewGinAccountRouter(cfg *config.Config) *gin.Engine {

	httpRouter := gin.Default()
	httpRouter.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[GIN-Custom-Log] %s - [%s] \"%s %s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006-01-02 - 15:04:05"),
			param.Request.Host,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Common Endpoint
	httpRouter.NoRoute(common.NotFound())

	reverseProxyController := controller.NewReverseProxyController(cfg, os.Getenv("BACKEND_ACCOUNT_HOSTNAME"))
	httpRouter.Any("/*proxyPath", middleware.TokenIntrospection(), reverseProxyController.Proxy)

	return httpRouter

}
