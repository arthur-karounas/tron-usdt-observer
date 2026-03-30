package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// --- Database Models ---

type TrackedWallet struct {
	gorm.Model
	Address       string `gorm:"uniqueIndex;not null"`
	LastTimestamp int64  `gorm:"default:0"`
}

type AllowedUser struct {
	gorm.Model
	TelegramID int64 `gorm:"uniqueIndex;not null"`
}

// Storage handles persistent data (PostgreSQL) and cache (Redis)
type Storage struct {
	db  *gorm.DB
	rdb *redis.Client
}

// New creates a new storage instance with auto-migration and connectivity checks
func New(pgDSN, redisAddr, redisPassword string) (*Storage, error) {
	// Initialize PostgreSQL with GORM
	db, err := gorm.Open(postgres.Open(pgDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Run database migrations
	if err := db.AutoMigrate(&TrackedWallet{}, &AllowedUser{}); err != nil {
		return nil, fmt.Errorf("failed to migrate postgres: %w", err)
	}

	// Setup Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	// Check Redis availability
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Storage{db: db, rdb: rdb}, nil
}

// --- Wallet Management ---

func (s *Storage) AddWallet(address string) error {
	wallet := TrackedWallet{Address: address}
	return s.db.FirstOrCreate(&wallet, TrackedWallet{Address: address}).Error
}

func (s *Storage) RemoveWallet(address string) error {
	return s.db.Where("address = ?", address).Delete(&TrackedWallet{}).Error
}

func (s *Storage) GetWallets() ([]TrackedWallet, error) {
	var wallets []TrackedWallet
	err := s.db.Find(&wallets).Error
	return wallets, err
}

func (s *Storage) UpdateWalletTimestamp(address string, timestamp int64) error {
	return s.db.Model(&TrackedWallet{}).Where("address = ?", address).Update("last_timestamp", timestamp).Error
}

// --- User Management ---

func (s *Storage) AddUser(id int64) error {
	user := AllowedUser{TelegramID: id}
	return s.db.FirstOrCreate(&user, AllowedUser{TelegramID: id}).Error
}

func (s *Storage) RemoveUser(id int64) error {
	return s.db.Where("telegram_id = ?", id).Delete(&AllowedUser{}).Error
}

func (s *Storage) GetUsers() ([]int64, error) {
	var users []AllowedUser
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	var ids []int64
	for _, u := range users {
		ids = append(ids, u.TelegramID)
	}
	return ids, nil
}

// --- Transaction Deduplication ---

// ProcessTransaction uses Redis SetNX to ensure a transaction is handled only once
func (s *Storage) ProcessTransaction(ctx context.Context, txID, address string) (bool, error) {
	key := fmt.Sprintf("seen_tx:%s", txID)
	// Transaction record expires after 7 days
	return s.rdb.SetNX(ctx, key, address, 7*24*time.Hour).Result()
}
