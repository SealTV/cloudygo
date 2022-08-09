package stability

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit[T any] func(ctx context.Context) (T, error)

func Breaker[T any](circuit Circuit[T], failureTreshhold uint) Circuit[T] {
	var consecutiveFailures int = 0
	var lastAttemt time.Time = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (T, error) {
		m.RLock()

		d := consecutiveFailures - int(failureTreshhold)
		if d >= 0 {
			shouldRetryAt := lastAttemt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return defaultVal[T](), errors.New("service unreachable")
			}
		}

		m.RUnlock()

		response, err := circuit(ctx)

		m.Lock()
		defer m.Unlock()

		if err != nil {
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0
		return response, nil
	}
}
