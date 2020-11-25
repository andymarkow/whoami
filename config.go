package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	ServerHost string `envconfig:"WEB_SERVER_HOST" default:"0.0.0.0"`
	ServerPort string `envconfig:"WEB_SERVER_PORT" default:"80"`
}

func getConfig() *config {
	cfg := &config{}

	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	return cfg
}
