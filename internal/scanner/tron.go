package scanner

import (
	"context"
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
