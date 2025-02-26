package config

import (
	"log"
	"os"
	"strconv"
)

// Config stores application configurations
type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresDB       string
	PostgresPort     int
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
}

// LoadConfig loads environment variables from .env file
func LoadConfig() *Config {

	postgresPort, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid POSTGRES_PORT value")
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value")
	}

	return &Config{
		PostgresUser:     getEnv("POSTGRES_USER", "user"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "password"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresDB:       getEnv("POSTGRES_DB", "dbname"),
		PostgresPort:     postgresPort,
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          redisDB,
	}
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
