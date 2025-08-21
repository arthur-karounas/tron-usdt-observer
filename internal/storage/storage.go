package storage

import (
	"gorm.io/gorm"
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
	db *gorm.DB
}
