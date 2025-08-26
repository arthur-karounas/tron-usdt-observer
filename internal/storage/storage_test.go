package storage

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessTransaction_Behavour(t *testing.T) {
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

	type testCase struct {
		name        string
		txId        string
		wallet      string
		expected    bool
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

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, gotNew)
			}
		})
	}
}
