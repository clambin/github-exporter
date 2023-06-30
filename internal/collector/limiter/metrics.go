package limiter

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var _ prometheus.Collector = &metrics{}

type metrics struct {
	inFlight          int
	maxInFlight       int
	inFlightMetric    *prometheus.Desc
	maxInFlightMetric *prometheus.Desc
	lock              sync.Mutex
}

func newMetrics(namespace, subsystem, application string) metrics {
	return metrics{
		inFlightMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "api_inflight"),
			"Number of requests in flight",
			nil,
			map[string]string{"application": application},
		),
		maxInFlightMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "api_max_inflight"),
			"Maximum number of requests in flight",
			nil,
			map[string]string{"application": application},
		),
	}
}

func (m *metrics) inc() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.inFlight++
	if m.inFlight > m.maxInFlight {
		m.maxInFlight = m.inFlight
	}
}

func (m *metrics) dec() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.inFlight--
}

func (m *metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.inFlightMetric
	ch <- m.maxInFlightMetric
}

func (m *metrics) Collect(ch chan<- prometheus.Metric) {
	m.lock.Lock()
	defer m.lock.Unlock()
	ch <- prometheus.MustNewConstMetric(m.inFlightMetric, prometheus.GaugeValue, float64(m.inFlight))
	ch <- prometheus.MustNewConstMetric(m.maxInFlightMetric, prometheus.GaugeValue, float64(m.maxInFlight))
}
