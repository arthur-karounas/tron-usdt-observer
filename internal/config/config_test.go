package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Behavour(t *testing.T) {
	type testCase struct {
		name     string
		input    map[string]string
		wantErr  bool
		expected *Config
	}

	testCases := []testCase{
		{
			name: "Success case with minimum required fields",
			input: map[string]string{
				"BOT_TOKEN": "12345:token",
				"ADMIN_ID":  "987654321",
			},
			wantErr: false,
			expected: &Config{
				BotToken:           "12345:token",
				AdminID:            987654321,
				PollInterval:       15,
				MaxConcurrentScans: 5,
				USDTContract:       "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			},
		},
		{
			name: "BOT_TOKEN is missing",
			input: map[string]string{
				"ADMIN_ID": "987654321",
			},
			wantErr: true,
		},
		{
			name: "ADMIN_ID is not a number",
			input: map[string]string{
				"ADMIN_ID": "not-a-number",
			},
			wantErr: true,
		},
		{
			name: "Success custom with Poll Interval and Max Scans",
			input: map[string]string{
				"BOT_TOKEN":            "12345:token",
				"ADMIN_ID":             "987654321",
				"POLL_INTERVAL":        "32",
				"MAX_CONCURRENT_SCANS": "10",
			},
			wantErr: false,
			expected: &Config{
				BotToken:           "12345:token",
				AdminID:            987654321,
				PollInterval:       32,
				MaxConcurrentScans: 10,
				USDTContract:       "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for key, value := range tt.input {
				os.Setenv(key, value)
			}

			cfg, err := Load()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.Equal(t, tt.expected.AdminID, cfg.AdminID)
				assert.Equal(t, tt.expected.BotToken, cfg.BotToken)
				assert.Equal(t, tt.expected.PollInterval, cfg.PollInterval)
				assert.Equal(t, tt.expected.MaxConcurrentScans, cfg.MaxConcurrentScans)
				assert.Equal(t, tt.expected.USDTContract, cfg.USDTContract)
			}
		})
	}
}
