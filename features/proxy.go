package features

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/drawbridge/config"
	"github.com/hexcraft-biz/drawbridge/controllers"
	"github.com/hexcraft-biz/drawbridge/middleware"
)

func LoadProxy(e *gin.Engine, cfg *config.Config) {
	c := controllers.NewProxyController(cfg)

	e.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[Proxy-Log] %s - [%s] \"%s %s %s %s %d %s \"%s\" %s\"\n",
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

	if cfg.ProxyAllowCORS == true {
		e.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           cfg.ProxyAllowCORSMaxAge,
		}))

		e.OPTIONS("/*options_support", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusNoContent)
			return
		})
	}

	/*
		// TODO if need.
		HEAD
		CONNECT
		TRACE
	*/

	e.GET(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.Dogmas(cfg),
		c.Proxy(),
	)

	e.POST(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.Dogmas(cfg),
		c.Proxy(),
	)

	e.PUT(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.Dogmas(cfg),
		c.Proxy(),
	)

	e.PATCH(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.Dogmas(cfg),
		c.Proxy(),
	)

	e.DELETE(
		"/*proxyPath",
		middleware.TokenIntrospection(cfg),
		middleware.Userinfo(cfg),
		middleware.Dogmas(cfg),
		c.Proxy(),
	)
}
