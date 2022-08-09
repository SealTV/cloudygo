package concurrency

import "sync"

func Funnel[T any](sources ...<-chan T) <-chan T {
	out := make(chan T)
	wg := sync.WaitGroup{}

	wg.Add(len(sources))

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ch := range sources {
		go func(in <-chan T) {
			defer wg.Done()
			for v := range in {
				out <- v
			}
		}(ch)
	}

	return out
}
