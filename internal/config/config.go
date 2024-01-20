package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	ServerHost         string
	ServerPort         string
	LogFormatter       string // Possible values: fmt, json.
	LogLevel           string // Possible values: error, warn, info, debug.
	AccessLogEnabled   bool
	AccessLogSkipPaths []string
	ReadTimeout        time.Duration
	ReadHeaderTimeout  time.Duration
	WriteTimeout       time.Duration
	TLSCrtFile         string
	TLSKeyFile         string
	TLSCAFile          string
}

// NewConfig creates a new Config object with default values.
func NewConfig() (*Config, error) {
	flag.Usage = func() {
		fmt.Printf("Whoami - Simple Go web server based on net/http library which returns information about web server and HTTP context.\n\n")
		fmt.Printf("Usage:\n")
		fmt.Printf("  whoami [flags]\n\n")
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}

	cfg := &Config{}

	var accessLogSkipPaths string
	var readTimeout, readHeaderTimeout, writeTimeout string

	flag.StringVar(&cfg.ServerHost, "host", getEnv("WHOAMI_HOST", "0.0.0.0"), "Web server host address")
	flag.StringVar(&cfg.ServerPort, "port", getEnv("WHOAMI_PORT", "8080"), "Web server port number")
	flag.StringVar(&cfg.LogFormatter, "log-formatter", getEnv("WHOAMI_LOG_FORMATTER", "json"), "Log formatter: 'fmt' or 'json'")
	flag.StringVar(&cfg.LogLevel, "log-level", getEnv("WHOAMI_LOG_LEVEL", "info"), "Log level: 'error', 'warn', 'error', 'debug'")
	flag.BoolVar(&cfg.AccessLogEnabled, "access-log", getEnv("WHOAMI_ACCESS_LOG", "false") == "true", "Enable access log")
	flag.StringVar(&accessLogSkipPaths, "access-log-skip-paths", getEnv("WHOAMI_ACCESS_LOG_SKIP_PATHS", ""), "Comma separated list of URL paths to skip in access log")
	flag.StringVar(&readTimeout, "read-timeout", getEnv("WHOAMI_READ_TIMEOUT", "0s"), "Web server read timeout")
	flag.StringVar(&readHeaderTimeout, "read-header-timeout", getEnv("WHOAMI_READ_HEADER_TIMEOUT", "0s"), "Web server read header timeout")
	flag.StringVar(&writeTimeout, "write-timeout", getEnv("WHOAMI_WRITE_TIMEOUT", "0s"), "Web server write timeout")
	flag.StringVar(&cfg.TLSCrtFile, "tls-crt", getEnv("WHOAMI_TLS_CRT_FILE", ""), "TLS certificate file")
	flag.StringVar(&cfg.TLSKeyFile, "tls-key", getEnv("WHOAMI_TLS_KEY_FILE", ""), "TLS private key file")
	flag.StringVar(&cfg.TLSCAFile, "tls-ca", getEnv("WHOAMI_TLS_CA_FILE", ""), "TLS CA certificate file for mTLS authentication")

	flag.Parse()

	if accessLogSkipPaths != "" {
		cfg.AccessLogSkipPaths = strings.Split(accessLogSkipPaths, ",")
	}

	var err error

	cfg.ReadTimeout, err = time.ParseDuration(readTimeout)
	if err != nil {
		return nil, fmt.Errorf("time.ParseDuration: %w", err)
	}

	cfg.ReadHeaderTimeout, err = time.ParseDuration(readHeaderTimeout)
	if err != nil {
		return nil, fmt.Errorf("time.ParseDuration: %w", err)
	}

	cfg.WriteTimeout, err = time.ParseDuration(writeTimeout)
	if err != nil {
		return nil, fmt.Errorf("time.ParseDuration: %w", err)
	}

	return cfg, nil
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
