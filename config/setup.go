package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"
)

// ================================================================
//
// ================================================================
type Config struct {
	*Env
	ProxyMappings *[]ProxyMapping
}

func Load() (*Config, error) {
	env, err := GetEnv()
	if err != nil {
		return nil, err
	}

	proxyMappings := []ProxyMapping{}

	getJson(env.ProxyMappginsFile, &proxyMappings)

	return &Config{Env: env, ProxyMappings: &proxyMappings}, nil
}

// ================================================================
//
// ================================================================
type Env struct {
	AppHostname        string
	AppPort            string
	GinMode            string
	Location           *time.Location
	ProxyMappginsFile  string
	ProxyAllowCORS     bool
	Oauth2AdminHost    string
	Oauth2PublicHost   string
	Oauth2ScopesHost   string
	Oauth2HeaderPrefix string
}

func GetEnv() (*Env, error) {
	var err error

	env := &Env{
		AppHostname: os.Getenv("APP_HOSTNAME"),
		AppPort:     os.Getenv("APP_PORT"),
		GinMode:     os.Getenv("GIN_MODE"),
	}

	if env.Location, err = time.LoadLocation(os.Getenv("TIMEZONE")); err != nil {
		return nil, err
	}

	if os.Getenv("PROXY_MAPPINGS_JSON_FILE_PATH") != "" {
		env.ProxyMappginsFile = os.Getenv("PROXY_MAPPINGS_JSON_FILE_PATH")
	} else {
		return nil, errors.New("Invalid environment variable : PROXY_MAPPINGS_JSON_FILE_PATH")
	}

	if os.Getenv("PROXY_ALLOW_CORS") != "" {
		if env.ProxyAllowCORS, err = strconv.ParseBool(os.Getenv("PROXY_ALLOW_CORS")); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid environment variable : PROXY_ALLOW_CORS")
	}

	if os.Getenv("OAUTH2_ADMIN_HOST") != "" {
		env.Oauth2AdminHost = os.Getenv("OAUTH2_ADMIN_HOST")
	} else {
		return nil, errors.New("Invalid environment variable : OAUTH2_ADMIN_HOST")
	}

	if os.Getenv("OAUTH2_PUBLIC_HOST") != "" {
		env.Oauth2PublicHost = os.Getenv("OAUTH2_PUBLIC_HOST")
	} else {
		return nil, errors.New("Invalid environment variable : OAUTH2_PUBLIC_HOST")
	}

	if os.Getenv("OAUTH2_SCOPES_HOST") != "" {
		env.Oauth2ScopesHost = os.Getenv("OAUTH2_SCOPES_HOST")
	} else {
		return nil, errors.New("Invalid environment variable : OAUTH2_SCOPES_HOST")
	}

	if os.Getenv("OAUTH2_HEADER_PREFIX") != "" {
		env.Oauth2HeaderPrefix = os.Getenv("OAUTH2_HEADER_PREFIX")
	} else {
		return nil, errors.New("Invalid environment variable : OAUTH2_HEADER_PREFIX")
	}

	return env, nil
}

// ================================================================
//
// ================================================================
type ProxyMapping struct {
	PublicHostname   string `json:"public-hostname"`
	InternalHostname string `json:"internal-hostname"`
}

func getJson(filePath string, target interface{}) error {
	data, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer data.Close()

	return json.NewDecoder(data).Decode(target)
}
