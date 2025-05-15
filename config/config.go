package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort    string
	Dsn           string
	RedisAddress  string
	RedisPassword string
}

func returnError(name string) error {
	fmt.Printf("%s not found", name)
	return fmt.Errorf("%s not found", name)
}

func SetupEnv() (AppConfig, error) {
	godotenv.Load()

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return AppConfig{}, returnError("http port")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		return AppConfig{}, returnError("database details")
	}

	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		return AppConfig{}, returnError("redis address")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	fmt.Printf("dsn is %s", dsn)

	return AppConfig{
		ServerPort:    httpPort,
		Dsn:           dsn,
		RedisAddress:  redisAddress,
		RedisPassword: redisPassword,
	}, nil
}
