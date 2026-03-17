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

type TrackedWallet struct {
	gorm.Model
	Address       string `gorm:"uniqueIndex;not null"`
	LastTimestamp int64  `gorm:"default:0"`
}

type AllowedUser struct {
	gorm.Model
	TelegramID int64 `gorm:"uniqueIndex;not null"`
}

type Storage struct {
	db  *gorm.DB
	rdb *redis.Client
}

func New(pgDSN, redisAddr, redisPassword string) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(pgDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err := db.AutoMigrate(&TrackedWallet{}, &AllowedUser{}); err != nil {
		return nil, fmt.Errorf("failed to migrate postgres: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Storage{db: db, rdb: rdb}, nil
}
