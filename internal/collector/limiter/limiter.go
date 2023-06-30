package limiter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/semaphore"
	"net/http"
)

var _ http.RoundTripper = &Limiter{}
var _ prometheus.Collector = &Limiter{}

type Limiter struct {
	parallel *semaphore.Weighted
	next     http.RoundTripper
	metrics
}

func New(maxParallel int64, next http.RoundTripper, namespace, subsystem, application string) *Limiter {
	if next == nil {
		next = http.DefaultTransport
	}
	return &Limiter{
		parallel: semaphore.NewWeighted(maxParallel),
		next:     next,
		metrics:  newMetrics(namespace, subsystem, application),
	}
}

func (l *Limiter) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := l.parallel.Acquire(request.Context(), 1); err != nil {
		return nil, fmt.Errorf("acquire semaphore: %w", err)
	}
	defer l.parallel.Release(1)

	l.metrics.inc()
	defer l.metrics.dec()

	return l.next.RoundTrip(request)
}
