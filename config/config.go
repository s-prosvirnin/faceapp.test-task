package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	CancelContextSleepDuration int    `envconfig:"CANCEL_CONTEXT_SLEEP_DURATION"`
	HttpHost                   string `envconfig:"HTTP_HOST"`
	HttpPort                   int    `envconfig:"HTTP_PORT"`

	DbHost     string `envconfig:"DB_HOST"`
	DbPort     int    `envconfig:"DB_PORT"`
	DbUser     string `envconfig:"DB_USER"`
	DbPassword string `envconfig:"DB_PASSWORD"`
	DbSchema   string `envconfig:"DB_SCHEMA"`
}

func InitConfig() (Config, error) {
	_ = godotenv.Load(".env")

	var cfg Config
	if err := envconfig.Process("TEST_TASK", &cfg); err != nil {
		return Config{}, errors.Wrap(err, "get config from env error")
	}

	return cfg, nil
}
