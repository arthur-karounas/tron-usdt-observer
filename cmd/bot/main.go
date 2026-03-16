package main

import (
	"github.com/arthur-karounas/tron-usdt-observer/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalf("Failed to load config: %v", err)
	}

	sugar.Infow("Configuration loaded successfully",
		"admin_id", cfg.AdminID,
		"poll_interval", cfg.PollInterval,
	)

	sugar.Info("Commit 1: Project structure and config initialized.")
}
