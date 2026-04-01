package scanner

import (
	"strings"
	"testing"
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
			if got != tt.expected {
				t.Errorf("ParseAmount(%s) = %f; want %f", tt.input, got, tt.expected)
			}
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
	if !strings.Contains(msg, "150.35 USDT") {
		t.Error("Message does not contain the correct amount")
	}

	// Verify wallet address truncation (security/privacy)
	if !strings.Contains(msg, "1234") || !strings.Contains(msg, "6543") {
		t.Error("Message does not contain wallet address tails")
	}

	// Verify transaction ID inclusion
	if !strings.Contains(msg, txID) {
		t.Error("Message does not contain the transaction link")
	}

	// Verify HTML styling for Telegram
	if !strings.Contains(msg, "<b>") || !strings.Contains(msg, "<code>") {
		t.Error("Message does not contain expected HTML tags")
	}
}
