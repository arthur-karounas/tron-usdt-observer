package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arthur-karounas/tron-usdt-observer/internal/bot"
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

	tgBot, err := bot.New(cfg, db, tronScanner, sugar)
	if err != nil {
		sugar.Fatalf("Bot initialization failed: %v", err)
	}

	tronScanner.SetNotifier(tgBot.SendNotification)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tronScanner.Start(ctx)
	}()

	go tgBot.Start()

	sugar.Info("System online. Press Ctrl+C to stop.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Shutting down gracefully...")
	cancel()
	tgBot.Stop()

	wg.Wait()
	sugar.Info("Shutdown complete.")
}
