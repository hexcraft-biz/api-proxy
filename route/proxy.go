package route

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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
		return fmt.Sprintf("[GIN-Proxy-Log] %s - [%s] \"%s %s %s %s %d %s \"%s\" %s\"\n",
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

	if cfg.Env.ProxyAllowCORS == true {
		httpRouter.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           1 * time.Hour,
		}))

		httpRouter.OPTIONS("/*options_support", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusNoContent)
			return
		})
	}

	reverseProxyController := controller.NewReverseProxyController(cfg, internalHostname)

	/*
		// TODO if need.
		HEAD
		CONNECT
		TRACE
	*/

	httpRouter.GET(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.VerifyScope(cfg, internalHostname),
		reverseProxyController.Proxy,
	)

	httpRouter.POST(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.VerifyScope(cfg, internalHostname),
		reverseProxyController.Proxy,
	)

	httpRouter.PUT(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.VerifyScope(cfg, internalHostname),
		reverseProxyController.Proxy,
	)

	httpRouter.PATCH(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.VerifyScope(cfg, internalHostname),
		reverseProxyController.Proxy,
	)

	httpRouter.DELETE(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.VerifyScope(cfg, internalHostname),
		reverseProxyController.Proxy,
	)

	return httpRouter

}
