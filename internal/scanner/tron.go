package scanner

import (
	"context"
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
