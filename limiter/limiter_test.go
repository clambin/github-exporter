package limiter

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimiter(t *testing.T) {
	var parallel, maxParallel atomic.Int64
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := parallel.Add(1)
		defer parallel.Add(-1)
		if current > maxParallel.Load() {
			maxParallel.Store(current)
		}
		time.Sleep(500 * time.Millisecond)
	}))

	const maxConcurrent = 5
	l := NewLimiter(maxConcurrent)
	httpClient := &http.Client{Transport: l.RoundTripper(http.DefaultTransport)}

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			_, err := httpClient.Get(ts.URL)
			require.NoError(t, err)
		})
	}
	wg.Wait()

	assert.Equal(t, int64(maxConcurrent), maxParallel.Load())
}
