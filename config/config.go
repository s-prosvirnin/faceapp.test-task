package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	CancelContextSleepDuration int    `envconfig:"CANCEL_CONTEXT_SLEEP_DURATION" default:"0"`
	HttpHost                   string `envconfig:"HTTP_HOST" default:"localhost"`
	HttpPort                   int    `envconfig:"HTTP_PORT" default:"8085"`

	DbHost     string `envconfig:"DB_HOST" default:"localhost"`
	DbPort     int    `envconfig:"DB_PORT" default:"54323"`
	DbUser     string `envconfig:"DB_USER" default:"postgres"`
	DbPassword string `envconfig:"DB_PASSWORD" default:"12345"`
	DbSchema   string `envconfig:"DB_SCHEMA" default:"postgres"`
}

func InitConfig() (Config, error) {
	_ = godotenv.Load(".env")

	var cfg Config
	if err := envconfig.Process("TEST_TASK", &cfg); err != nil {
		return Config{}, errors.Wrap(err, "get config from env error")
	}

	return cfg, nil
}
