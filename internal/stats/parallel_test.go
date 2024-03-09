package stats

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	var p parallel[int]

	var errEven = errors.New("even")

	for i := range 1000 {
		p.Do(func() (int, error) {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			var err error
			if i%2 == 0 {
				err = errEven
			}
			return i, err
		})
	}

	numbers, err := p.Results()
	assert.ErrorIs(t, err, errEven)
	assert.Len(t, numbers, 500)
}
