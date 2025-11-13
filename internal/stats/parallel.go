package stats

import (
	"errors"
	"sync"
)

type parallel[T any] struct {
	err    error
	result []T
	wg     sync.WaitGroup
	lock   sync.RWMutex
}

func (p *parallel[T]) Do(f func() (T, error)) {
	// we don't need to limit concurrent calls, as the calling app uses a limiting round tripper
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		val, err := f()
		p.lock.Lock()
		defer p.lock.Unlock()
		p.err = errors.Join(p.err, err)
		if err == nil {
			p.result = append(p.result, val)
		}
	}()
}

func (p *parallel[T]) Results() ([]T, error) {
	p.wg.Wait()
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.result, p.err
}
