package route

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/api-proxy/common"
	"github.com/hexcraft-biz/api-proxy/config"
	"github.com/hexcraft-biz/api-proxy/controller"
	"github.com/hexcraft-biz/api-proxy/middleware"
)

// all the routes are defined here
func NewGinProxyRouter(cfg *config.Config, internalHostname string) *gin.Engine {

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

	reverseProxyController := controller.NewReverseProxyController(cfg, internalHostname)
	httpRouter.Any("/*proxyPath", middleware.TokenIntrospection(cfg), reverseProxyController.Proxy)

	return httpRouter

}
