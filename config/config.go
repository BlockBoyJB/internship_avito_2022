package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP  HTTP
	Log   Log
	PG    PG
	Redis Redis
	JWT   JWT
	Kafka Kafka
}

type (
	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}
	Log struct {
		Level  string `env-required:"true" env:"LOG_LEVEL"`
		Output string `env-required:"true" env:"LOG_OUTPUT"`
	}
	PG struct {
		MaxPoolSize int    `env-required:"true" env:"PG_MAX_POOL_SIZE"`
		Url         string `env-required:"true" env:"PG_URL"`
	}
	Redis struct {
		Url      string `env-required:"true" env:"REDIS_URL"`
		Password string `env:"REDIS_PASSWORD"`
	}
	JWT struct {
		PrivateKey string `env-required:"true" env:"JWT_PRIVATE_KEY"`
		PublicKey  string `env-required:"true" env:"JWT_PUBLIC_KEY"`
	}
	Kafka struct {
		Url string `env-required:"true" env:"KAFKA_URL"`
	}
)

func NewConfig() (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadEnv(c); err != nil {
		return nil, fmt.Errorf("error reading config env: %w", err)
	}
	return c, nil
}
