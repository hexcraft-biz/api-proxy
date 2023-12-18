package main

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/drawbridge/config"
	"github.com/hexcraft-biz/drawbridge/features"
)

type App struct {
	mainHandler http.Handler
	HostSwitch  HostSwitch
}

type HostSwitch map[string]http.Handler

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestHost := r.Host
	if host, _, err := net.SplitHostPort(r.Host); err == nil {
		requestHost = host
	}

	if handler := a.HostSwitch[requestHost]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		a.mainHandler.ServeHTTP(w, r)
	}
}

func main() {
	cfg, err := config.Load()
	MustNot(err)

	app := App{
		mainHandler: GetProxyEngine(cfg),
		HostSwitch: HostSwitch{
			cfg.AppHost: GetCommonEngine(cfg),
		},
	}

	http.ListenAndServe(":"+cfg.AppPort, app)
}

func GetProxyEngine(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{cfg.TrustProxy})

	features.LoadProxy(r, cfg)

	return r
}

func GetCommonEngine(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{cfg.TrustProxy})

	features.LoadCommon(r, cfg)

	return r
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
