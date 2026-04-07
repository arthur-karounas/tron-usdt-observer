package obs

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all prometheus metric collectors
type Metrics struct {
	// Scanner Metrics
	ScansTotal        *prometheus.CounterVec
	ScanErrorsTotal   *prometheus.CounterVec
	ScanDuration      *prometheus.HistogramVec
	TransactionsFound *prometheus.CounterVec
	TrackedWallets    prometheus.Gauge
	ScannerStatus     prometheus.Gauge // 1 for running, 0 for stopped

	// Bot Metrics
	NotificationsSent *prometheus.CounterVec
	BotCommandsTotal  *prometheus.CounterVec
}

// NewMetrics initializes and registers application metrics in the global registry.
// Use this for production code.
func NewMetrics() *Metrics {
	return &Metrics{
		ScansTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "tron_observer_scans_total",
			Help: "The total number of wallet scan attempts",
		}, []string{"status"}),

		ScanErrorsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "tron_observer_scan_errors_total",
			Help: "The total number of errors encountered during scanning",
		}, []string{"type"}),

		ScanDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "tron_observer_scan_duration_seconds",
			Help:    "Time taken to scan a single wallet",
			Buckets: prometheus.DefBuckets,
		}, []string{"wallet"}),

		TransactionsFound: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "tron_observer_transactions_found_total",
			Help: "Number of new USDT transactions detected",
		}, []string{"wallet"}),

		TrackedWallets: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "tron_observer_tracked_wallets_count",
			Help: "Current number of wallets being monitored",
		}),

		ScannerStatus: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "tron_observer_scanner_running_status",
			Help: "Operational status of the scanner (1=Running, 0=Stopped)",
		}),

		NotificationsSent: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "tron_observer_notifications_sent_total",
			Help: "Total telegram notifications sent",
		}, []string{"status"}),

		BotCommandsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "tron_observer_bot_commands_total",
			Help: "Total number of bot commands received",
		}, []string{"command"}),
	}
}

// NewTestMetrics creates metrics without registering them in the global Prometheus registry.
// This prevents "duplicate metrics registration" errors during tests.
func NewTestMetrics() *Metrics {
	return &Metrics{
		ScansTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "test_scans_total",
		}, []string{"status"}),

		ScanErrorsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "test_scan_errors_total",
		}, []string{"type"}),

		ScanDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "test_scan_duration_seconds",
		}, []string{"wallet"}),

		TransactionsFound: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "test_transactions_found_total",
		}, []string{"wallet"}),

		TrackedWallets: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "test_tracked_wallets_count",
		}),

		ScannerStatus: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "test_scanner_status",
		}),

		NotificationsSent: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "test_notifications_sent_total",
		}, []string{"status"}),

		BotCommandsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "test_bot_commands_total",
		}, []string{"command"}),
	}
}
