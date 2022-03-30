package main

import (
	"net/http"
	"os"

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
	hostSwitch[os.Getenv("APP_HOSTNAME")+":"+os.Getenv("APP_PORT")] = route.NewGinMainRouter(cfg)
	hostSwitch[os.Getenv("PUBLIC_ACCOUNT_HOSTNAME")+":"+os.Getenv("APP_PORT")] = route.NewGinAccountRouter(cfg)

	http.ListenAndServe(":"+cfg.AppPort, hostSwitch)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
