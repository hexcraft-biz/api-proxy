package main

import (
	"fmt"
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
	fmt.Println("request from", r.Host)
	if handler := a.HostSwitch[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		a.mainHandler.ServeHTTP(w, r)
	}
}

func main() {
	cfg, err := config.Load()
	MustNot(err)

	hostSwitch := make(HostSwitch)

	app := App{mainHandler: route.NewGinMainRouter(cfg), HostSwitch: hostSwitch}

	for _, pm := range *cfg.ProxyMappings {
		fmt.Println("proxy-setting", pm.PublicHostname, pm.InternalHostname)
		app.HostSwitch[pm.PublicHostname+":"+cfg.Env.AppPort] = route.NewGinProxyRouter(cfg, pm.InternalHostname)
	}

	http.ListenAndServe(":"+cfg.AppPort, app)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
