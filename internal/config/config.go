package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultPort     = "8080"
	defaultRedisURL = "redis://localhost:6380/0"
)

type Config struct {
	Port     string
	RedisURL string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = defaultRedisURL
	}

	return Config{
		Port:     port,
		RedisURL: redisURL,
	}, nil
}
