// Package config is for config
package config

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	RedisAddr     string
	JWTPublicKey  ed25519.PublicKey
	JWTPrivateKey ed25519.PrivateKey
}
type rawConfig struct {
	ServerPort    string `env:"SERVER_PORT" envDefault:"8080"`
	RedisAddr     string `env:"REDIS_ADDR,required"`
	JWTPublicKey  string `env:"JWT_PUBLIC_KEY,required"`
	JWTPrivateKey string `env:"JWT_PRIVATE_KEY,required"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found. Loading from the environment")
	}

	rawCfg := &rawConfig{}

	if err := env.Parse(rawCfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse config: %v", err)
	}

	publicKey, err := base64.StdEncoding.DecodeString(rawCfg.JWTPublicKey)
	if err != nil {
		return nil, fmt.Errorf("config: error decoding jwt public key: %v", err)
	}

	privateKey, err := base64.StdEncoding.DecodeString(rawCfg.JWTPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("config: error decoding jwt public key: %v", err)
	}

	cfg := &Config{
		ServerPort:    rawCfg.ServerPort,
		RedisAddr:     rawCfg.RedisAddr,
		JWTPublicKey:  ed25519.PublicKey(publicKey),
		JWTPrivateKey: ed25519.PrivateKey(privateKey),
	}

	return cfg, nil
}
