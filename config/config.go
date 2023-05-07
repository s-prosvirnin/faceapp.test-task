package config

import "time"

// упростил - параметры захардкожены вместо пробрасывания через env
type Config struct {
	CancelContextSleepDuration time.Duration
	HttpHost                   string
	HttpPort                   int

	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbSchema   string
}

func InitConfig() Config {
	return Config{
		CancelContextSleepDuration: 2 * time.Second,
		HttpHost:                   "localhost",
		HttpPort:                   8085,
		DbHost:                     "localhost",
		DbPort:                     54323,
		DbUser:                     "postgres",
		DbPassword:                 "12345",
		DbSchema:                   "postgres",
	}
}
