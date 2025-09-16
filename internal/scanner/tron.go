package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/arthur-karounas/tron-usdt-observer/internal/config"
	"github.com/arthur-karounas/tron-usdt-observer/internal/storage"
	"go.uber.org/zap"
)

// Store defines database operations required by the scanner
type Store interface {
	GetWallets() ([]storage.TrackedWallet, error)
	UpdateWalletTimestamp(address string, timestamp int64) error
	ProcessTransaction(ctx context.Context, txID, address string) (bool, error)
}

// Scanner monitors the Tron blockchain for USDT transfers
type Scanner struct {
	cfg        *config.Config
	db         Store
	logger     *zap.SugaredLogger
	client     *http.Client
	isRunning  bool
	runMutex   sync.RWMutex
	notifyFunc func(msg string)
}

// New initializes a new blockchain scanner
func New(cfg *config.Config, db Store, logger *zap.SugaredLogger) *Scanner {
	return &Scanner{
		cfg:       cfg,
		db:        db,
		logger:    logger,
		client:    &http.Client{Timeout: 5 * time.Second},
		isRunning: true,
	}
}

// --- Data Formatting Helpers ---

// ParseAmount converts raw Sun (10^-6) string to float64 USDT
func ParseAmount(raw string) float64 {
	var amount float64
	fmt.Sscanf(raw, "%f", &amount)
	return amount / 1_000_000
}

// FormatNotification prepares an HTML message for Telegram
func FormatNotification(address, from string, amount float64, timestamp int64, txID string) string {
	loc, _ := time.LoadLocation("Europe/Warsaw")
	t := time.UnixMilli(timestamp).In(loc)
	timeStr := t.Format("02 January 2006 (15:04)")

	addrTail := address
	if len(address) > 4 {
		addrTail = address[len(address)-4:]
	}
	fromTail := from
	if len(from) > 4 {
		fromTail = from[len(from)-4:]
	}

	return fmt.Sprintf("📥 <b>Incoming USDT Transaction</b>\n\n"+
		"Wallet: <code>...%s</code>\n"+
		"Sender: <code>...%s</code>\n"+
		"Amount: <b>%.2f USDT</b>\n"+
		"Time: %s\n\n"+
		"<a href='https://tronscan.org/#/transaction/%s'>View Transaction</a>",
		addrTail, fromTail, amount, timeStr, txID)
}

// --- Lifecycle & Status ---

func (s *Scanner) SetNotifier(f func(string)) {
	s.notifyFunc = f
}

func (s *Scanner) SetRunning(state bool) {
	s.runMutex.Lock()
	defer s.runMutex.Unlock()
	s.isRunning = state
}

func (s *Scanner) IsRunning() bool {
	s.runMutex.RLock()
	defer s.runMutex.RUnlock()
	return s.isRunning
}

// Start kicks off the background polling loop
func (s *Scanner) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.cfg.PollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Scanner stopped by context")
			return
		case <-ticker.C:
			if !s.IsRunning() {
				continue
			}
			s.processWallets(ctx)
		}
	}
}

// --- Core Logic ---

// processWallets iterates through tracked addresses using concurrency
func (s *Scanner) processWallets(ctx context.Context) {
	wallets, err := s.db.GetWallets()
	if err != nil || len(wallets) == 0 {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, s.cfg.MaxConcurrentScans) // Concurrency limit

	for _, w := range wallets {
		if ctx.Err() != nil {
			return
		}

		wg.Add(1)
		go func(wallet storage.TrackedWallet) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			if ctx.Err() != nil {
				return
			}

			// Initialize timestamp for new wallets (last 24h)
			if wallet.LastTimestamp == 0 {
				wallet.LastTimestamp = time.Now().Add(-24 * time.Hour).UnixMilli()
				s.db.UpdateWalletTimestamp(wallet.Address, wallet.LastTimestamp)
			}

			s.fetchAndProcessTransactions(ctx, wallet)
		}(w)
	}

	wg.Wait()
}

// fetchAndProcessTransactions calls TronGrid API and handles new transfers
func (s *Scanner) fetchAndProcessTransactions(ctx context.Context, w storage.TrackedWallet) {
	url := fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions/trc20?limit=50&contract_address=%s&only_to=true&min_timestamp=%d",
		w.Address, s.cfg.USDTContract, w.LastTimestamp)

	var resp *http.Response
	var err error

	// Retry logic for API rate limits and network errors
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		req.Header.Set("Accept", "application/json")
		if s.cfg.TronAPIKey != "" {
			req.Header.Set("TRON-PRO-API-KEY", s.cfg.TronAPIKey)
		}

		resp, err = s.client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		// Error handling and backoff
		if resp != nil {
			resp.Body.Close()
			if resp.StatusCode == 403 || resp.StatusCode == 429 {
				s.logger.Warnf("Rate limit exceeded for %s (Attempt %d/%d)", w.Address, attempt+1, maxRetries+1)
			} else {
				s.logger.Warnf("API returned status %d (Attempt %d/%d)", resp.StatusCode, attempt+1, maxRetries+1)
			}
		} else {
			s.logger.Warnf("HTTP error: %v (Attempt %d/%d)", err, attempt+1, maxRetries+1)
		}

		if attempt == maxRetries {
			s.logger.Errorf("Failed to fetch transactions for %s after %d attempts", w.Address, maxRetries+1)
			return
		}

		delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
		time.Sleep(delay)
	}

	defer resp.Body.Close()

	// Parse JSON response structure
	var result struct {
		Success bool `json:"success"`
		Data    []struct {
			TransactionID  string `json:"transaction_id"`
			From           string `json:"from"`
			Value          string `json:"value"`
			BlockTimestamp int64  `json:"block_timestamp"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Success {
		s.logger.Errorf("Failed to decode JSON: %v", err)
		return
	}

	var maxTimestamp int64 = w.LastTimestamp

	// Process transactions from oldest to newest
	for i := len(result.Data) - 1; i >= 0; i-- {
		tx := result.Data[i]

		// Check if transaction was already processed (Redis cache)
		isNew, err := s.db.ProcessTransaction(ctx, tx.TransactionID, w.Address)
		if err != nil {
			s.logger.Errorf("Redis error: %v", err)
			continue
		}
		if !isNew {
			continue
		}

		amount := ParseAmount(tx.Value)
		msg := FormatNotification(w.Address, tx.From, amount, tx.BlockTimestamp, tx.TransactionID)

		// Async notification delivery
		if s.notifyFunc != nil {
			go s.notifyFunc(msg)
		}

		if tx.BlockTimestamp > maxTimestamp {
			maxTimestamp = tx.BlockTimestamp
		}
	}

	// Update bookmark for next poll
	if maxTimestamp > w.LastTimestamp {
		s.db.UpdateWalletTimestamp(w.Address, maxTimestamp)
	}

	// Small pause to prevent hitting rate limits too fast
	time.Sleep(250 * time.Millisecond)
}
