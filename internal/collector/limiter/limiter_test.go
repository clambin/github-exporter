package limiter_test

import (
	"bytes"
	"context"
	"github.com/clambin/github-exporter/internal/collector/limiter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestLimiter_RoundTrip(t *testing.T) {
	reg := prometheus.NewPedanticRegistry()

	const maxParallel = 10
	s := stubbedServer{delay: 10 * time.Millisecond}
	r := limiter.New(maxParallel, &s, "foo", "bar", "snafu")
	c := http.Client{Transport: r}

	reg.MustRegister(r)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := c.Get("/")
			require.NoError(t, err)
		}()
	}
	wg.Wait()
	assert.LessOrEqual(t, s.maxInFlight, maxParallel)

	assert.NoError(t, testutil.GatherAndCompare(reg, bytes.NewBufferString(`
# HELP foo_bar_api_max_inflight Maximum number of requests in flight
# TYPE foo_bar_api_max_inflight gauge
foo_bar_api_max_inflight{application="snafu"} 10
`), "foo_bar_max_inflight"))
}

func TestLimiter_RoundTrip_Exceeded(t *testing.T) {
	s := stubbedServer{delay: time.Second}
	r := limiter.New(1, &s, "foo", "bar", "snafu")
	c := http.Client{Transport: r}

	go func() {
		_, _ = c.Get("/")
	}()

	// wait for the first request to reach the server
	assert.Eventually(t, func() bool {
		s.lock.Lock()
		defer s.lock.Unlock()

		return s.called > 0
	}, time.Second, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	_, err := c.Do(req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

var _ http.RoundTripper = &stubbedServer{}

type stubbedServer struct {
	delay       time.Duration
	called      int
	inFlight    int
	maxInFlight int
	lock        sync.Mutex
}

func (s *stubbedServer) RoundTrip(_ *http.Request) (*http.Response, error) {
	s.inc()
	defer s.dec()

	time.Sleep(s.delay)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`hello`)),
	}, nil
}

func (s *stubbedServer) inc() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.called++
	s.inFlight++
	if s.inFlight > s.maxInFlight {
		s.maxInFlight = s.inFlight
	}
}

func (s *stubbedServer) dec() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.inFlight--
}
