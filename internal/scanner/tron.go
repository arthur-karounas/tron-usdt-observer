package scanner

import (
	"context"

	"github.com/arthur-karounas/tron-usdt-observer/internal/storage"
)

type Store interface {
	GetWallets() ([]storage.TrackedWallet, error)
	UpdateWalletTimestamp(address string, timestamp int64) error
	ProcessTransaction(ctx context.Context, txID, address string) (bool, error)
}
