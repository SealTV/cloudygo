package concurrency

func Split[T any](in <-chan T, n int) [](<-chan T) {
	dest := make([](<-chan T), 0, n)

	for i := 0; i < n; i++ {
		ch := make(chan T)
		dest = append(dest, ch)

		go func() {
			defer close(ch)

			for val := range in {
				ch <- val
			}
		}()
	}

	return dest
}
