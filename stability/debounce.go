package stability

import (
	"context"
	"sync"
	"time"

	"github.com/SealTV/cloudygo/helper"
)

func DebounceFirst[T any](circuit Circuit[T], d time.Duration) Circuit[T] {
	var threshold time.Time
	var result T
	var err error
	var m sync.Mutex

	return func(ctx context.Context) (T, error) {
		m.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if !time.Now().Before(threshold) {
			result, err = circuit(ctx)
		}

		return result, err
	}
}

func DebounceLast[T any](circuit Circuit[T], d time.Duration) Circuit[T] {
	var (
		threshold time.Time
		ticker    *time.Ticker
		result    T
		err       error
		once      sync.Once
		m         sync.Mutex
	)

	return func(ctx context.Context) (T, error) {
		m.Lock()
		defer m.Unlock()

		threshold = time.Now().Add(d)

		once.Do(func() {
			ticker = time.NewTicker(100 * time.Millisecond)

			go func() {
				defer func() {
					m.Lock()
					defer m.Unlock()

					ticker.Stop()
					once = sync.Once{}
				}()

				for {
					select {
					case <-ticker.C:
						m.Lock()
						if time.Now().Before(threshold) {
							result, err = circuit(ctx)
							m.Unlock()
							return
						}
					case <-ctx.Done():
						m.Lock()
						result, err = helper.DefaultVal[T](), ctx.Err()
						m.Unlock()
						return
					}
				}
			}()
		})

		return result, err
	}
}
