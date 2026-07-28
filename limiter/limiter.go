package limiter

import (
	"net/http"

	"golang.org/x/sync/semaphore"
)

type Limiter struct {
	sema *semaphore.Weighted
}

func NewLimiter(max int64) *Limiter {
	return &Limiter{
		sema: semaphore.NewWeighted(max),
	}
}

func (l *Limiter) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		if err := l.sema.Acquire(request.Context(), 1); err != nil {
			return nil, err
		}
		defer l.sema.Release(1)
		return next.RoundTrip(request)
	})
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (r roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return r(request)
}
