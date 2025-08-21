package bot

import (
	"github.com/arthur-karounas/tron-usdt-observer/internal/storage"
)

type Store interface {
	GetUsers() ([]int64, error)
	GetWallets() ([]storage.TrackedWallet, error)
	AddWallet(address string) error
	RemoveWallet(address string) error
	AddUser(id int64) error
	RemoveUser(id int64) error
}

type ScannerController interface {
	SetRunning(state bool)
	IsRunning() bool
}
