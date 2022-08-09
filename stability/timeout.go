package stability

import "context"

type SlowFunction[T any] func(T) (T, error)

type WithConext[T any] func(context.Context, T) (T, error)

func Timeout[T any](f SlowFunction[T]) WithConext[T] {
	return func(ctx context.Context, arg T) (T, error) {
		chres := make(chan T)
		cherr := make(chan error)

		go func() {
			res, err := f(arg)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		case <-ctx.Done():
			return defaultVal[T](), ctx.Err()
		}
	}
}
