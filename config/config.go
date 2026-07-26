package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerPort     string
	Environment    string
	AllowedOrigins []string

	DatabaseURL     string
	DatabaseMaxConn int
	DatabaseMaxIdle int

	RedisURL      string
	RedisPassword string
	RedisDB       int

	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	RateLimitRPS   int
	RateLimitBurst int
}

func Load() *Config {
	return &Config{
		ServerPort:     getEnv("PORT", "8080"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		AllowedOrigins: getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),

		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/topup_game?sslmode=disable"),
		DatabaseMaxConn: getEnvInt("DATABASE_MAX_CONN", 25),
		DatabaseMaxIdle: getEnvInt("DATABASE_MAX_IDLE", 5),

		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		JWTAccessExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 24*time.Hour),

		RateLimitRPS:   getEnvInt("RATE_LIMIT_RPS", 10),
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 100),
	}
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err != nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
