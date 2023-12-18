package main

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/drawbridge/config"
	"github.com/hexcraft-biz/drawbridge/features"
)

type App struct {
	proxyHandler http.Handler
	hostSwitch   hostSwitch
}

type hostSwitch map[string]http.Handler

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the host from the request
	requestHost := r.Host

	// If the host contains a port, remove it to get just the host
	if host, _, err := net.SplitHostPort(r.Host); err == nil {
		requestHost = host
	}

	// Check if the host is an IP address
	if ip := net.ParseIP(requestHost); ip != nil {
		// If it's an IP address, use the default handler ("main")
		a.hostSwitch["main"].ServeHTTP(w, r)
	} else if handler := a.hostSwitch[requestHost]; handler != nil {
		// If there's a specific handler for the host, use it
		handler.ServeHTTP(w, r)
	} else {
		// If no specific handler is found, use the default proxy handler
		a.proxyHandler.ServeHTTP(w, r)
	}
}

func main() {
	cfg, err := config.Load()
	MustNot(err)

	commenEngine := GetCommonEngine(cfg)

	app := App{
		proxyHandler: GetProxyEngine(cfg),
		hostSwitch: hostSwitch{
			"main":      commenEngine,
			cfg.AppHost: commenEngine,
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
