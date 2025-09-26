package storage

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Test deduplication logic using in-memory Redis mock
func TestProcessTransaction_Behavour(t *testing.T) {
	// Setup miniredis (mock)
	mr, err := miniredis.Run()
	require.NoError(t, err, "failed to run miniredis (mock redis)")
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	s := &Storage{
		rdb: rdb,
	}

	ctx := context.Background()

	// Define test cases for transaction processing
	type testCase struct {
		name        string
		txId        string
		wallet      string
		expected    bool // Should be marked as 'new'
		expectedErr bool
	}

	testCases := []testCase{
		{name: "First time seeing transaction", txId: "123", wallet: "abc", expected: true, expectedErr: false},
		{name: "Second time seeing same transaction", txId: "123", wallet: "abc", expected: false, expectedErr: false},
		{name: "Same transaction but another wallet", txId: "123", wallet: "def", expected: false, expectedErr: false},
		{name: "Same wallet but another transaction", txId: "456", wallet: "abc", expected: true, expectedErr: false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gotNew, err := s.ProcessTransaction(ctx, tt.txId, tt.wallet)

			// Validate results
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, gotNew)
			}
		})
	}
}

// Test wallet management operations using SQLite in-memory
func TestWalletOperations(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db.AutoMigrate(&TrackedWallet{})

	s := &Storage{db: db}
	addr := "  TAddress123  "

	// Test AddWallet
	err = s.AddWallet(addr)
	assert.NoError(t, err)

	// Test GetWallets
	wallets, err := s.GetWallets()
	assert.NoError(t, err)
	assert.Len(t, wallets, 1)
	assert.Equal(t, "TAddress123", wallets[0].Address)

	// Test UpdateWalletTimestamp
	newTS := int64(123456789)
	err = s.UpdateWalletTimestamp("TAddress123", newTS)
	assert.NoError(t, err)

	wallets, _ = s.GetWallets()
	assert.Equal(t, newTS, wallets[0].LastTimestamp)

	// Test RemoveWallet
	err = s.RemoveWallet("TAddress123")
	assert.NoError(t, err)
	wallets, _ = s.GetWallets()
	assert.Empty(t, wallets)
}

// Test user management operations using SQLite in-memory
func TestUserOperations(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	db.AutoMigrate(&AllowedUser{})

	s := &Storage{db: db}
	userID := int64(998877)

	// Test AddUser
	err = s.AddUser(userID)
	assert.NoError(t, err)

	// Test GetUsers
	users, err := s.GetUsers()
	assert.NoError(t, err)
	assert.Contains(t, users, userID)

	// Test RemoveUser
	err = s.RemoveUser(userID)
	assert.NoError(t, err)
	users, _ = s.GetUsers()
	assert.NotContains(t, users, userID)
}
