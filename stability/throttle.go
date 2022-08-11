package stability

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/SealTV/cloudygo/helper"
)

type Effector[T any] func(ctx context.Context) (T, error)

func Throttle[T any](e Effector[T], max, refill uint, d time.Duration) Effector[T] {
	var (
		tokens uint = max
		once   sync.Once
	)

	return func(ctx context.Context) (T, error) {
		if err := ctx.Err(); err != nil {
			return helper.DefaultVal[T](), err
		}

		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}

						tokens = t
					}
				}
			}()
		})

		if tokens <= 0 {
			return helper.DefaultVal[T](), errors.New("too many calls")
		}
		tokens--
		return e(ctx)
	}
}
