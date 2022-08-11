package concurrency

import (
	"context"
	"sync"
	"time"

	"github.com/SealTV/cloudygo/helper"
)

type Future[T any] interface {
	Result() (T, error)
}

type InnerFuture[T any] struct {
	once sync.Once
	wg   sync.WaitGroup

	res     T
	err     error
	rChan   <-chan T
	errChan <-chan error
}

func (f *InnerFuture[T]) Result() (T, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()

		f.res = <-f.rChan
		f.err = <-f.errChan
	})

	f.wg.Wait()
	return f.res, f.err
}

func SlowFunction[T any](ctx context.Context) Future[T] {
	rChan := make(chan T)
	errChan := make(chan error)

	go func() {
		select {
		case <-time.After(2 * time.Second):
			rChan <- helper.DefaultVal[T]()
			errChan <- nil
		case <-ctx.Done():
			rChan <- helper.DefaultVal[T]()
			errChan <- ctx.Err()
		}
	}()

	return &InnerFuture[T]{rChan: rChan, errChan: errChan}
}
