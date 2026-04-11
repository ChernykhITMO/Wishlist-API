package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBConfig   DBConfig
	HTTPConfig HTTPConfig
	JWTSecret  string
}

type DBConfig struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnIdleTime time.Duration
	MaxConnLifetime time.Duration
}

type HTTPConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	const op = "config.Load"

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("%s: DB_URL is required", op)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("%s: JWT_SECRET is required", op)
	}

	return &Config{
		DBConfig: DBConfig{
			URL:             dbURL,
			MaxConns:        getEnvInt32("DB_MAX_CONNS", 10),
			MinConns:        getEnvInt32("DB_MIN_CONNS", 2),
			MaxConnIdleTime: getEnvDuration("DB_MAX_CONN_IDLE_TIME", 5*time.Minute),
			MaxConnLifetime: getEnvDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute),
		},
		HTTPConfig: HTTPConfig{
			Addr:            getEnvString("HTTP_ADDR", ":8080"),
			ReadTimeout:     getEnvDuration("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getEnvDuration("HTTP_WRITE_TIMEOUT", 5*time.Second),
			ShutdownTimeout: getEnvDuration("HTTP_SHUTDOWN_TIMEOUT", 5*time.Second),
		},
		JWTSecret: jwtSecret,
	}, nil
}

func getEnvString(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	return val
}

func getEnvInt32(key string, def int32) int32 {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return int32(i)
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return def
	}

	return d
}
