package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Port        int
	PostgresDSN string

	// DB fields used to build PostgresDSN when POSTGRES_DSN is not set
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads configuration from environment variables.
// If a .env file exists in the current directory, it is loaded first.
// PostgresDSN is set from POSTGRES_DSN, or built from DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE.
func Load() (*Config, error) {
	_ = godotenv.Load()

	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}

	dbPort := 5432
	if p := os.Getenv("DB_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			dbPort = v
		}
	}

	dbSSLMode := "disable"
	if s := os.Getenv("DB_SSLMODE"); s != "" {
		dbSSLMode = s
	}

	cfg := &Config{
		Port:        port,
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  dbSSLMode,
	}

	if cfg.PostgresDSN == "" && (cfg.DBHost != "" || cfg.DBUser != "" || cfg.DBPassword != "" || cfg.DBName != "") {
		cfg.PostgresDSN = buildDSN(cfg)
	}

	return cfg, nil
}

func buildDSN(c *Config) string {
	host := c.DBHost
	if host == "" {
		host = "localhost"
	}
	user := c.DBUser
	if user == "" {
		user = "postgres"
	}
	password := c.DBPassword
	dbname := c.DBName
	if dbname == "" {
		dbname = "behavix"
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, c.DBPort, dbname, c.DBSSLMode)
}
