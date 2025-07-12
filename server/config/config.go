// Package config is for config
package config

import "os"

type Config struct {
	Port string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		return fallback
	}
}
