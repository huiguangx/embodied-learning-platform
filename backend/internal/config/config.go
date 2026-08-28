package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port int
}

func Load() Config {
	return Config{
		Host: env("HOST", "0.0.0.0"),
		Port: envInt("PORT", 8080),
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
