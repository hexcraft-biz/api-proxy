package controllers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/drawbridge/config"
	"github.com/hexcraft-biz/drawbridge/constants"
	"github.com/hexcraft-biz/feature"
	"github.com/hexcraft-biz/her"
)

type ProxyController struct {
	Config *config.Config
}

func NewProxyController(cfg *config.Config) *ProxyController {
	return &ProxyController{cfg}
}

func (r *ProxyController) Proxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.MustGet(constants.MiddlewareKeyProxyRoute).(*feature.Route)
		if remote, err := url.Parse(route.RootUrl); err != nil {
			c.AbortWithStatusJSON(her.NewErrorWithMessage(http.StatusInternalServerError, "drawbridge: "+err.Error(), nil).HttpR())
		} else {
			c.Request.Header.Del("Authorization")

			proxy := httputil.NewSingleHostReverseProxy(remote)
			// TODO: full proxy route:
			// 	Method
			//	RootUrl
			//	Feature
			//	Path
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
}
