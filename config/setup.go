package config

import (
	"os"
	"time"
)

//================================================================
//
//================================================================
type Config struct {
	*Env
}

func Load() (*Config, error) {
	env, err := GetEnv()
	if err != nil {
		return nil, err
	}

	return &Config{Env: env}, nil
}

//================================================================
//
//================================================================
type Env struct {
	AppHostname string
	AppPort     string
	GinMode     string
	Location    *time.Location
}

func GetEnv() (*Env, error) {
	var err error

	var loc *time.Location
	if loc, err = time.LoadLocation(os.Getenv("TIMEZONE")); err != nil {
		return nil, err
	}

	return &Env{
		AppHostname: os.Getenv("APP_HOSTNAME"),
		AppPort:     os.Getenv("APP_PORT"),
		GinMode:     os.Getenv("GIN_MODE"),
		Location:    loc,
	}, nil
}
