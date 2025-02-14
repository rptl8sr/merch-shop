package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	Prod = "prod"
	Dev  = "dev"
)

type Config struct {
	// App
	AppLogLevel slog.Level `envconfig:"APP_LOG_LEVEL" default:"info"`
	AppPort     int        `envconfig:"APP_PORT" default:"8080"`

	// PostgreSQL
	PGHost     string `envconfig:"PG_HOST" required:"true"`
	PGPort     int    `envconfig:"PG_PORT" required:"true"`
	PGUser     string `envconfig:"PG_USER" required:"true"`
	PGPassword string `envconfig:"PG_PASSWORD" required:"true"`
	PGDatabase string `envconfig:"PG_DATABASE" required:"true"`
	PGSSLMode  string `envconfig:"PG_SSL_MODE" default:"disable"`

	// JWT
	JWTSecret string `envconfig:"JWT_SECRET" required:"true"`
}

func Must() *Config {
	env := os.Getenv("ENV")

	switch env {
	case Prod:
	case Dev:
		if err := godotenv.Load(".env", fmt.Sprintf(".env.%s", Dev)); err != nil {
			panic(fmt.Sprintf("Error loading '.env' or '.env.dev' files: %v", err))
		}
	default:
		panic(fmt.Sprintf("ENV environment variable must be set to '%s' or '%s'", Prod, Dev))
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic(fmt.Sprintf("Error processing environment variables: %v", err))
	}

	return &cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.PGHost, c.PGPort, c.PGUser, c.PGPassword, c.PGDatabase, c.PGSSLMode)
}
