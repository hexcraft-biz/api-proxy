package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	app "github.com/hexcraft-biz/envmod-app"
	"github.com/hexcraft-biz/feature"
)

// ================================================================
//
// ================================================================
type Config struct {
	*app.App
	*feature.Dogmas
	ProxyAllowCORS         bool
	ProxyAllowCORSMaxAge   time.Duration
	OAuth2AdminHost        string
	OAuth2PublicHost       string
	OAuth2HeaderInfix      string
	ContextKeyTargetPrefix string
}

func Load() (*Config, error) {
	emApp, err := app.New()
	if err != nil {
		return nil, err
	}

	emDogmas, err := feature.NewDogmas(emApp.AppRootUrl)
	if err != nil {
		return nil, err
	}

	config := &Config{
		App:                    emApp,
		Dogmas:                 emDogmas,
		OAuth2HeaderInfix:      os.Getenv("OAUTH2_HEADER_INFIX"),
		ContextKeyTargetPrefix: "api-proxy-target-",
	}

	if os.Getenv("PROXY_ALLOW_CORS") != "" {
		if config.ProxyAllowCORS, err = strconv.ParseBool(os.Getenv("PROXY_ALLOW_CORS")); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid environment variable : PROXY_ALLOW_CORS")
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("PROXY_ALLOW_CORS_MAX_AGE")); err != nil {
		return nil, err
	} else if exist == true {
		config.ProxyAllowCORSMaxAge = time.Duration(value) * time.Second
	}

	if os.Getenv("OAUTH2_ADMIN_HOST") != "" {
		config.OAuth2AdminHost = os.Getenv("OAUTH2_ADMIN_HOST")
	} else {
		return nil, errors.New("Invalid environment variable : OAUTH2_ADMIN_HOST")
	}

	if os.Getenv("OAUTH2_PUBLIC_HOST") != "" {
		config.OAuth2PublicHost = os.Getenv("OAUTH2_PUBLIC_HOST")
	} else {
		return nil, errors.New("Invalid environment variable : OAUTH2_PUBLIC_HOST")
	}

	if os.Getenv("OAUTH2_HEADER_INFIX") != "" {
		config.OAuth2HeaderInfix = os.Getenv("OAUTH2_HEADER_INFIX")
	} else {
		return nil, errors.New("Invalid environment variable : OAUTH2_HEADER_INFIX")
	}

	return config, nil
}

func FetchOptIntEnv(envStr string) (value int, exist bool, err error) {
	if envStr != "" {
		exist = true
		if intVal, err := strconv.Atoi(envStr); err != nil {
			return value, exist, err
		} else {
			value = intVal
		}
	} else {
		exist = false
	}

	return value, exist, nil
}
