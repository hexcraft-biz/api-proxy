package main

import (
	"net/http"

	"github.com/hexcraft-biz/api-proxy/config"
	"github.com/hexcraft-biz/api-proxy/route"
)

type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
}

func main() {
	cfg, err := config.Load()
	MustNot(err)

	hostSwitch := make(HostSwitch)
	hostSwitch[cfg.Env.AppHostname+":"+cfg.Env.AppPort] = route.NewGinMainRouter(cfg)

	for _, pm := range *cfg.ProxyMappings {
		hostSwitch[pm.PublicHostname+":"+cfg.Env.AppPort] = route.NewGinProxyRouter(cfg, pm.InternalHostname)
	}

	http.ListenAndServe(":"+cfg.AppPort, hostSwitch)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
