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

type Store interface {
	GetWallets() ([]storage.TrackedWallet, error)
	UpdateWalletTimestamp(address string, timestamp int64) error
	ProcessTransaction(ctx context.Context, txID, address string) (bool, error)
}

type Scanner struct {
	cfg        *config.Config
	db         Store
	logger     *zap.SugaredLogger
	client     *http.Client
	isRunning  bool
	runMutex   sync.RWMutex
	notifyFunc func(msg string)
}

func New(cfg *config.Config, db Store, logger *zap.SugaredLogger) *Scanner {
	return &Scanner{
		cfg:       cfg,
		db:        db,
		logger:    logger,
		client:    &http.Client{Timeout: 5 * time.Second},
		isRunning: true,
	}
}

func ParseAmount(raw string) float64 {
	var amount float64
	fmt.Sscanf(raw, "%f", &amount)
	return amount / 1_000_000
}

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

func (s *Scanner) fetchAndProcessTransactions(ctx context.Context, w storage.TrackedWallet) {
	url := fmt.Sprintf("https://api.trongrid.io/v1/accounts/%s/transactions/trc20?limit=50&contract_address=%s&only_to=true&min_timestamp=%d",
		w.Address, s.cfg.USDTContract, w.LastTimestamp)

	var resp *http.Response
	var err error

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

		delay := baseDelay * time.Duration(1<<attempt)
		time.Sleep(delay)
	}

	defer resp.Body.Close()

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

	for i := len(result.Data) - 1; i >= 0; i-- {
		tx := result.Data[i]

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

		if s.notifyFunc != nil {
			go s.notifyFunc(msg)
		}

		if tx.BlockTimestamp > maxTimestamp {
			maxTimestamp = tx.BlockTimestamp
		}
	}

	if maxTimestamp > w.LastTimestamp {
		s.db.UpdateWalletTimestamp(w.Address, maxTimestamp)
	}

	time.Sleep(250 * time.Millisecond)
}
