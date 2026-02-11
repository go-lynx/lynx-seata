package seata

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SeataMetrics defines Seata-related monitoring metrics.
type SeataMetrics struct {
	transactionsTotal   *prometheus.CounterVec
	transactionsActive  prometheus.Gauge
	transactionDuration *prometheus.HistogramVec
	branchTransactions  prometheus.Gauge
	healthCheckTotal    *prometheus.CounterVec
	healthCheckDuration *prometheus.HistogramVec
}

var (
	metrics     *SeataMetrics
	metricsOnce sync.Once
)

// getMetrics returns the metrics instance (lazy init).
func (t *TxSeataClient) getMetrics() *SeataMetrics {
	if !t.IsEnabled() {
		return nil
	}
	metricsOnce.Do(func() {
		metrics = newSeataMetrics()
	})
	return metrics
}

func newSeataMetrics() *SeataMetrics {
	return &SeataMetrics{
		transactionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "lynx",
				Subsystem: "seata",
				Name:      "transactions_total",
				Help:      "Total number of Seata transactions",
			},
			[]string{"status"},
		),
		transactionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "lynx",
				Subsystem: "seata",
				Name:      "transactions_active",
				Help:      "Number of active Seata transactions",
			},
		),
		transactionDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "lynx",
				Subsystem: "seata",
				Name:      "transaction_duration_seconds",
				Help:      "Duration of Seata transactions",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"status"},
		),
		branchTransactions: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "lynx",
				Subsystem: "seata",
				Name:      "branch_transactions",
				Help:      "Number of Seata branch transactions",
			},
		),
		healthCheckTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "lynx",
				Subsystem: "seata",
				Name:      "health_check_total",
				Help:      "Total number of Seata health checks",
			},
			[]string{"status"},
		),
		healthCheckDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "lynx",
				Subsystem: "seata",
				Name:      "health_check_duration_seconds",
				Help:      "Duration of Seata health checks",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"status"},
		),
	}
}

// RecordTransaction records a transaction completion.
func (m *SeataMetrics) RecordTransaction(status string) {
	if m != nil && m.transactionsTotal != nil {
		m.transactionsTotal.WithLabelValues(status).Inc()
	}
}

// RecordTransactionDuration records transaction duration.
func (m *SeataMetrics) RecordTransactionDuration(status string, duration float64) {
	if m != nil && m.transactionDuration != nil {
		m.transactionDuration.WithLabelValues(status).Observe(duration)
	}
}

// SetActiveTransactions sets the number of active transactions.
func (m *SeataMetrics) SetActiveTransactions(count float64) {
	if m != nil && m.transactionsActive != nil {
		m.transactionsActive.Set(count)
	}
}

// SetBranchTransactions sets the number of branch transactions.
func (m *SeataMetrics) SetBranchTransactions(count float64) {
	if m != nil && m.branchTransactions != nil {
		m.branchTransactions.Set(count)
	}
}

// RecordHealthCheck records a health check.
func (m *SeataMetrics) RecordHealthCheck(status string) {
	if m != nil && m.healthCheckTotal != nil {
		m.healthCheckTotal.WithLabelValues(status).Inc()
	}
}

// GetMetrics returns the Seata metrics instance for external use.
func GetMetrics() *SeataMetrics {
	return metrics
}
