package controller

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/karmaksana-io/api-proxy/config"
)

type ReverseProxyController struct {
	Config      *config.Config
	ProxyTarget string
}

func NewReverseProxyController(cfg *config.Config, proxyTarget string) ReverseProxyController {
	return ReverseProxyController{cfg, proxyTarget}
}

func (r ReverseProxyController) Proxy(ctx *gin.Context) {

	remote, err := url.Parse(r.ProxyTarget)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = ctx.Param("proxyPath")
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
