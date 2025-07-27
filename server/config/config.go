// Package config is for config
package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	privatePemPath = "secrets/private.pem"
	publicPemPath  = "secrets/public.pem"
)

type Config struct {
	ServerPort    string
	RedisAddr     string
	JWTPublicKey  *rsa.PublicKey
	JWTPrivateKey *rsa.PrivateKey
}
type rawConfig struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	RedisAddr  string `env:"REDIS_ADDR,required"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found. Loading from the environment")
	}

	rawCfg := &rawConfig{}

	if err := env.Parse(rawCfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse config: %v", err)
	}

	publicKey, err := loadPublicKey(publicPemPath)
	if err != nil {
		return nil, fmt.Errorf("config: error loading public key, %v", err)
	}

	privateKey, err := loadPrivateKey(privatePemPath)
	if err != nil {
		return nil, fmt.Errorf("config: error loading private key, %v", err)
	}

	cfg := &Config{
		ServerPort:    rawCfg.ServerPort,
		RedisAddr:     rawCfg.RedisAddr,
		JWTPublicKey:  publicKey,
		JWTPrivateKey: privateKey,
	}

	return cfg, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)

	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := privInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid private key")
	}

	return key, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key")
	}

	return key, nil
}
