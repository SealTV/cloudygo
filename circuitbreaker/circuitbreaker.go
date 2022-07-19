package circuitbraker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, failureTreshhold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttemt time.Time = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()

		d := consecutiveFailures - int(failureTreshhold)
		if d >= 0 {
			shouldRetryAt := lastAttemt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreacheable")
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
