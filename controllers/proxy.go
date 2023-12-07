package controllers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/api-proxy/config"
)

type ProxyController struct {
	Config *config.Config
}

func NewProxyController(cfg *config.Config) *ProxyController {
	return &ProxyController{cfg}
}

func (r *ProxyController) Proxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		targetHost := c.GetString(r.Config.ContextKeyTargetPrefix + "rootUrl")

		remote, err := url.Parse(targetHost)
		if err != nil {
			panic(err)
		}

		c.Request.Header.Del("Authorization")

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Param("proxyPath")
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
