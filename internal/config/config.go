package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort  string
	Environment string
	DatabaseURL string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret            string
	JWTExpirationMinutes int

	BookingTimeoutMinutes int
}

func LoadConfig() (*Config, error) {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_MINUTES", "60"))
	bookingTimeout, _ := strconv.Atoi(getEnv("BOOKING_TIMEOUT_MINUTES", "15"))

	config := &Config{
		ServerPort:            getEnv("SERVER_PORT", "8080"),
		Environment:           getEnv("ENVIRONMENT", "development"),
		DatabaseURL:           getEnv("DATABASE_URL", ""),
		RedisAddr:             getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", ""),
		RedisDB:               redisDB,
		JWTSecret:             getEnv("JWT_SECRET", "AAA"),
		JWTExpirationMinutes:  jwtExpiration,
		BookingTimeoutMinutes: bookingTimeout,
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}