package main

import "time"

// упростил - параметры захардкожены вместо пробрасывания через env
type config struct {
	cancelContextSleepDuration time.Duration
	httpHost                   string
	httpPort                   int

	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbSchema   string
}

func initConfig() config {
	return config{
		cancelContextSleepDuration: 2 * time.Second,
		httpHost:                   "localhost",
		httpPort:                   8085,
		dbHost:                     "localhost",
		dbPort:                     54323,
		dbUser:                     "postgres",
		dbPassword:                 "12345",
		dbSchema:                   "postgres",
	}
}
