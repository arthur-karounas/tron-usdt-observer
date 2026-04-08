package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arthur-karounas/tron-usdt-observer/internal/bot"
	"github.com/arthur-karounas/tron-usdt-observer/internal/config"
	"github.com/arthur-karounas/tron-usdt-observer/internal/obs"
	"github.com/arthur-karounas/tron-usdt-observer/internal/scanner"
	"github.com/arthur-karounas/tron-usdt-observer/internal/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Load environment configuration
	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalf("Failed to load config: %v", err)
	}

	// Initialize storage with retry logic
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

	// Initialize Observability
	metrics := obs.NewMetrics()
	go func() {
		sugar.Info("Starting metrics server on :2112/metrics")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			sugar.Errorf("Metrics server failed: %v", err)
		}
	}()

	// Setup initial data and context
	db.AddUser(cfg.AdminID)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize core components (Scanner & Bot)
	tronScanner := scanner.New(cfg, db, sugar, metrics)
	tgBot, err := bot.New(cfg, db, tronScanner, sugar, metrics)
	if err != nil {
		sugar.Fatalf("Bot initialization failed: %v", err)
	}

	// Link scanner to bot notifications
	tronScanner.SetNotifier(tgBot.SendNotification)

	// Run scanner in background
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tronScanner.Start(ctx)
	}()

	// Start Telegram bot
	go tgBot.Start()
	sugar.Info("System online. Press Ctrl+C to stop.")

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sugar.Info("Shutting down gracefully...")
	cancel()
	tgBot.Stop()

	wg.Wait()
	sugar.Info("Shutdown complete.")
}
