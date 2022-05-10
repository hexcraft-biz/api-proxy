package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/hexcraft-biz/api-proxy/config"
	"github.com/hexcraft-biz/api-proxy/route"
)

type App struct {
	mainHandler http.Handler
	HostSwitch  HostSwitch
}

type HostSwitch map[string]http.Handler

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[TEST-LOG] request from", r.Host)

	requestHost := r.Host
	if host, _, err := net.SplitHostPort(r.Host); err == nil {
		requestHost = host
	}

	if handler := a.HostSwitch[requestHost]; handler != nil {
		fmt.Println("[TEST-LOG] Go proxy route", r.Host, requestHost)
		handler.ServeHTTP(w, r)
	} else {
		fmt.Println("[TEST-LOG] Go main route")
		a.mainHandler.ServeHTTP(w, r)
	}
}

func main() {
	cfg, err := config.Load()
	MustNot(err)

	hostSwitch := make(HostSwitch)

	app := App{mainHandler: route.NewGinMainRouter(cfg), HostSwitch: hostSwitch}

	fmt.Println("[TEST-LOG] ENV : ", cfg.Env)

	for _, pm := range *cfg.ProxyMappings {
		fmt.Println("[TEST-LOG] proxy-setting : ", pm.PublicHostname, pm.InternalHostname)
		app.HostSwitch[pm.PublicHostname] = route.NewGinProxyRouter(cfg, pm.InternalHostname)
	}

	http.ListenAndServe(":"+cfg.AppPort, app)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
