package config

import (
	"os"
	"time"
)

const (
	DefaultLogLevel        = "info"
	DefaultGracefulTimeout = 5 * time.Second
)

type Config struct {
	LogLevel        string
	BinanceAPIKey   string
	BinanceSecret   string
	GracefulTimeout time.Duration
}

func NewConfig() (*Config, error) {
	// Default value for GracefulTimeout
	gracefulTimeout := DefaultGracefulTimeout

	// Override if there is an environment variable
	if gt := os.Getenv("GRACEFUL_TIMEOUT"); gt != "" {
		var err error
		gracefulTimeout, err = time.ParseDuration(gt)
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		LogLevel:        getEnv("LOG_LEVEL", DefaultLogLevel),
		BinanceAPIKey:   os.Getenv("BINANCE_API_KEY"),
		BinanceSecret:   os.Getenv("BINANCE_SECRET"),
		GracefulTimeout: gracefulTimeout,
	}, nil
}

// Helper function to get environment variables with a fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
