package main

import (
	"context"
	"sync"
	"time"

	"github.com/arthur-karounas/tron-usdt-observer/internal/config"
	"github.com/arthur-karounas/tron-usdt-observer/internal/scanner"
	"github.com/arthur-karounas/tron-usdt-observer/internal/storage"
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

	var db *storage.Storage
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		db, err = storage.New(cfg.PgDSN, cfg.RedisAddr, cfg.RedisPassword)
		if err == nil {
			break
		}

		if i == maxRetries {
			sugar.Fatalf("Storage initialization failed after %d attempts: %v", maxRetries, err)
		}

		sugar.Warnf("Failed to connect to storage, retrying in 3s... (Attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(3 * time.Second)
	}

	db.AddUser(cfg.AdminID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tronScanner := scanner.New(cfg, db, sugar)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tronScanner.Start(ctx)
	}()

	sugar.Info("Scanner logic online. Press Ctrl+C to exit (not graceful yet).")

	select {}
}
