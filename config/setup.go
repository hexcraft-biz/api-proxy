package config

import (
	"encoding/json"
	"os"
	"time"
)

//================================================================
//
//================================================================
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

//================================================================
//
//================================================================
type Env struct {
	AppHostname        string
	AppPort            string
	GinMode            string
	Location           *time.Location
	ProxyMappginsFile  string
	Oauth2AdminHost    string
	Oauth2PublicHost   string
	Oauth2ScopesHost   string
	Oauth2HeaderPrefix string
}

func GetEnv() (*Env, error) {
	var err error

	var loc *time.Location
	if loc, err = time.LoadLocation(os.Getenv("TIMEZONE")); err != nil {
		return nil, err
	}

	return &Env{
		AppHostname:        os.Getenv("APP_HOSTNAME"),
		AppPort:            os.Getenv("APP_PORT"),
		GinMode:            os.Getenv("GIN_MODE"),
		ProxyMappginsFile:  os.Getenv("PROXY_MAPPINGS_JSON_FILE_PATH"),
		Location:           loc,
		Oauth2AdminHost:    os.Getenv("OAUTH2_ADMIN_HOST"),
		Oauth2PublicHost:   os.Getenv("OAUTH2_PUBLIC_HOST"),
		Oauth2ScopesHost:   os.Getenv("OAUTH2_SCOPES_HOST"),
		Oauth2HeaderPrefix: os.Getenv("OAUTH2_HEADER_PREFIX"),
	}, nil
}

//================================================================
//
//================================================================
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
