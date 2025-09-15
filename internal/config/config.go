package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all environment-based application settings
type Config struct {
	BotToken           string
	AdminID            int64
	TronAPIKey         string
	PollInterval       int
	MaxConcurrentScans int
	USDTContract       string
	PgDSN              string
	RedisAddr          string
	RedisPassword      string
}

// Load fetches configuration from environment variables
func Load() (*Config, error) {
	// Mandatory: Telegram Bot credentials
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}

	adminIDStr := os.Getenv("ADMIN_ID")
	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil || adminID == 0 {
		return nil, fmt.Errorf("valid ADMIN_ID is required")
	}

	// Optional: Scanner behavior (with defaults)
	pollInterval := 15
	if pi := os.Getenv("POLL_INTERVAL"); pi != "" {
		if val, err := strconv.Atoi(pi); err == nil {
			pollInterval = val
		}
	}

	maxScans := 5
	if ms := os.Getenv("MAX_CONCURRENT_SCANS"); ms != "" {
		if val, err := strconv.Atoi(ms); err == nil && val > 0 {
			maxScans = val
		}
	}

	// Blockchain constants
	contract := os.Getenv("USDT_CONTRACT")
	if contract == "" {
		contract = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	}

	// Database connections (PostgreSQL & Redis)
	pgDSN := os.Getenv("PG_DSN")
	if pgDSN == "" {
		pgDSN = "host=localhost user=postgres password=postgres dbname=tron_db port=5432 sslmode=disable TimeZone=Europe/Warsaw"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return &Config{
		BotToken:           token,
		AdminID:            adminID,
		TronAPIKey:         os.Getenv("TRON_API_KEY"),
		PollInterval:       pollInterval,
		MaxConcurrentScans: maxScans,
		USDTContract:       contract,
		PgDSN:              pgDSN,
		RedisAddr:          redisAddr,
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
	}, nil
}
