package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ServerHost       string
	ServerPort       string
	LogFormatter     string // Possible values: fmt, json.
	LogLevel         string // Possible values: error, warn, info, debug.
	AccessLogEnabled bool
}

// NewConfig creates a new Config object with default values.
//
// Return:
// - *Config: a pointer to the newly created Config object.
func NewConfig() *Config {
	flag.Usage = func() {
		fmt.Printf("Whoami - simple web server for development and testing purposes.\n\n")
		fmt.Printf("Usage:\n")
		fmt.Printf("  whoami [flags]\n\n")
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}

	cfg := &Config{}

	flag.StringVar(&cfg.ServerHost, "host", getEnv("WHOAMI_HOST", "0.0.0.0"), "Web server host address")
	flag.StringVar(&cfg.ServerPort, "port", getEnv("WHOAMI_PORT", "8080"), "Web server port number")
	flag.StringVar(&cfg.LogFormatter, "log-formatter", getEnv("WHOAMI_LOG_FORMATTER", "json"), "Log formatter: 'fmt' or 'json'")
	flag.StringVar(&cfg.LogLevel, "log-level", getEnv("WHOAMI_LOG_LEVEL", "info"), "Log level: 'error', 'warn', 'error', 'debug'")
	flag.BoolVar(&cfg.AccessLogEnabled, "access-log", getEnv("WHOAMI_ACCESS_LOG", "false") == "true", "Enable access log")

	flag.Parse()

	return cfg
}

// getEnv gets the value of an environment variable specified by the key.
//
// It takes two parameters:
//   - key: a string representing the name of the environment variable.
//   - fallback: a string representing the fallback value to be returned if the environment variable is empty.
//
// It returns a string representing the value of the environment variable if it is set, otherwise it returns the fallback value.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
