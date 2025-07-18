// Package config is for config
package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	RedisAddr  string `env:"REDIS_ADDR,required"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found. Loading from the environment")
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
