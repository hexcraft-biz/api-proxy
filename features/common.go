package features

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/api-proxy/config"
	"github.com/hexcraft-biz/api-proxy/controllers"
)

func LoadCommon(e *gin.Engine, cfg *config.Config) {
	c := controllers.NewCommon(cfg)
	e.NoRoute(c.NotFound())

	e.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[Main-Log] %s - [%s] \"%s %s %s %s %d %s \"%s\" %s\"\n",
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

	f := e.Group("/healthcheck/v1")
	f.GET("/ping", c.Ping())
}
