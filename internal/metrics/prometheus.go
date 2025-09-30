package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	internetUp      prometheus.Gauge
	corporateUp     prometheus.Gauge
	checkDuration   prometheus.Gauge
	checksTotal     *prometheus.CounterVec
	checksSuccess   *prometheus.CounterVec
	checksFailed    *prometheus.CounterVec
}

func NewPrometheusMetrics(port int) (*PrometheusMetrics, error) {
	metrics := &PrometheusMetrics{
		internetUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "nexa_internet_up",
			Help: "Internet reachable (1/0)",
		}),
		corporateUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "nexa_corporate_up", 
			Help: "Corporate reachable (1/0)",
		}),
		checkDuration: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "nexa_check_duration_seconds",
			Help: "Duration of the last check in seconds",
		}),
		checksTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "nexa_checks_total",
			Help: "Total number of checks performed",
		}, []string{"type"}),
		checksSuccess: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "nexa_checks_success_total",
			Help: "Total number of successful checks",
		}, []string{"type"}),
		checksFailed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "nexa_checks_failed_total", 
			Help: "Total number of failed checks",
		}, []string{"type"}),
	}

	prometheus.MustRegister(
		metrics.internetUp,
		metrics.corporateUp,
		metrics.checkDuration,
		metrics.checksTotal,
		metrics.checksSuccess,
		metrics.checksFailed,
	)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Starting Prometheus metrics server on %s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Printf("Failed to start Prometheus server: %v\n", err)
		}
	}()

	return metrics, nil
}

func (m *PrometheusMetrics) UpdateInternetStatus(up bool) {
	if up {
		m.internetUp.Set(1)
	} else {
		m.internetUp.Set(0)
	}
}

func (m *PrometheusMetrics) UpdateCorporateStatus(up bool) {
	if up {
		m.corporateUp.Set(1)
	} else {
		m.corporateUp.Set(0)
	}
}

func (m *PrometheusMetrics) UpdateCheckDuration(duration float64) {
	m.checkDuration.Set(duration)
}

func (m *PrometheusMetrics) UpdateCheckSummary(stats struct {
	TotalChecks    int
	Successful     int
	Failed         int
	ExternalChecks int
	CorporateChecks int
}) {
	m.checksTotal.WithLabelValues("external").Add(float64(stats.ExternalChecks))
	m.checksSuccess.WithLabelValues("external").Add(float64(stats.Successful))
	m.checksFailed.WithLabelValues("external").Add(float64(stats.Failed))

	m.checksTotal.WithLabelValues("corporate").Add(float64(stats.CorporateChecks))
	m.checksSuccess.WithLabelValues("corporate").Add(float64(stats.Successful))
	m.checksFailed.WithLabelValues("corporate").Add(float64(stats.Failed))
}