package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Version        string
	ServerAddr     string   `envconfig:"WHOAMI_SERVER_ADDR" default:":8080"`
	LogLevel       string   `envconfig:"WHOAMI_LOG_LEVEL" default:"info"`
	LogURLExcludes []string `envconfig:"WHOAMI_LOG_URL_EXCLUDES" default:"/favicon.ico,/healthz,/metrics"`
}

func New(version string) *Config {
	cfg := &Config{}
	cfg.Version = version

	if err := envconfig.Process("", cfg); err != nil {
		panic(err)
	}

	return cfg
}
