package scanner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/arthur-karounas/tron-usdt-observer/internal/config"
	"github.com/arthur-karounas/tron-usdt-observer/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Test USDT amount parsing (Sun to USDT conversion)
func TestParseAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"1 USDT", "1000000", 1.0},
		{"0.5 USDT", "500000", 0.5},
		{"120.75 USDT", "120750000", 120.75},
		{"Zero", "0", 0.0},
		{"Small amount", "1", 0.000001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseAmount(tt.input)
			assert.Equal(t, tt.expected, got, "ParseAmount(%s)", tt.input)
		})
	}
}

// Test HTML notification message building
func TestFormatNotification(t *testing.T) {
	const (
		wallet = "TWDJL7p1234"
		sender = "TABC9876543"
		amount = 150.35
		ts     = 1710000000000
		txID   = "tx_hash_example_123"
	)

	msg := FormatNotification(wallet, sender, amount, ts, txID)

	// Verify amount formatting
	assert.Contains(t, msg, "150.35 USDT")

	// Verify wallet address truncation (security/privacy)
	assert.Contains(t, msg, "1234")
	assert.Contains(t, msg, "6543")

	// Verify transaction ID inclusion
	assert.Contains(t, msg, txID)

	// Verify HTML styling for Telegram
	assert.Contains(t, msg, "<b>")
	assert.Contains(t, msg, "<code>")
}

type mockStore struct {
	sync.Mutex
	wallets []storage.TrackedWallet
	updated map[string]int64
	txs     map[string]bool
}

func (m *mockStore) GetWallets() ([]storage.TrackedWallet, error) {
	m.Lock()
	defer m.Unlock()
	return m.wallets, nil
}

func (m *mockStore) UpdateWalletTimestamp(addr string, ts int64) error {
	m.Lock()
	defer m.Unlock()
	if m.updated == nil {
		m.updated = make(map[string]int64)
	}
	m.updated[addr] = ts
	return nil
}

func (m *mockStore) ProcessTransaction(ctx context.Context, txID, addr string) (bool, error) {
	m.Lock()
	defer m.Unlock()
	if m.txs == nil {
		m.txs = make(map[string]bool)
	}
	isNew := !m.txs[txID]
	m.txs[txID] = true
	return isNew, nil
}

// Test processing wallets and fetching transactions successfully
func TestScanner_ProcessWallets_Success(t *testing.T) {
	cfg := &config.Config{MaxConcurrentScans: 1, USDTContract: "TR7", TronAPIKey: "key"}
	db := &mockStore{
		wallets: []storage.TrackedWallet{
			{Address: "TNewWallet", LastTimestamp: 0},
			{Address: "TOldWallet", LastTimestamp: 100},
		},
	}
	s := New(cfg, db, zap.NewNop().Sugar())

	s.retryDelay = 1 * time.Microsecond
	s.apiPause = 1 * time.Microsecond

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "key", r.Header.Get("TRON-PRO-API-KEY"))

		parts := strings.Split(r.URL.Path, "/")
		var addr string
		if len(parts) > 3 {
			addr = parts[3]
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"data": [
				{"transaction_id": "tx123_` + addr + `", "from": "Sender", "value": "15000000", "block_timestamp": 200}
			],
			"meta": {"at": 200}
		}`))
	}))
	defer ts.Close()

	s.apiBaseURL = ts.URL

	var notifications []string
	var notifyMu sync.Mutex
	s.SetNotifier(func(msg string) {
		notifyMu.Lock()
		defer notifyMu.Unlock()
		notifications = append(notifications, msg)
	})

	s.processWallets(context.Background())

	time.Sleep(10 * time.Millisecond)

	db.Lock()
	defer db.Unlock()

	assert.NotZero(t, db.updated["TNewWallet"])
	assert.Equal(t, int64(200), db.updated["TOldWallet"])
	assert.True(t, db.txs["tx123_TNewWallet"])
	assert.True(t, db.txs["tx123_TOldWallet"])

	notifyMu.Lock()
	assert.Len(t, notifications, 2)
	notifyMu.Unlock()
}

// Test scanner state management and notifier setup
func TestScanner_Lifecycle(t *testing.T) {
	s := New(&config.Config{}, &mockStore{}, zap.NewNop().Sugar())

	assert.True(t, s.IsRunning())

	s.SetRunning(false)
	assert.False(t, s.IsRunning())

	notified := false
	s.SetNotifier(func(msg string) { notified = true })
	assert.NotNil(t, s.notifyFunc)

	s.notifyFunc("test message")
	assert.True(t, notified)
}
