package concurrency

import (
	"sort"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplit(t *testing.T) {
	in := make(chan int)

	chansCount := 3
	chans := Split(in, chansCount)

	if len(chans) != chansCount {
		t.Fatalf("invalid channels count, expect: %d, got: %d", chansCount, len(chans))
	}

	mx := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(chans))

	result := []int{}
	for _, c := range chans {
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				mx.Lock()
				result = append(result, v)
				mx.Unlock()
			}
		}(c)
	}

	values := []int{1, 2, 3, 4, 5, 6}
	for _, v := range values {
		in <- v
	}
	close(in)

	wg.Wait()
	sort.Ints(result)
	if diff := cmp.Diff(values, result); diff != "" {
		t.Errorf("Funnel mismatch: {-want; +got}: %s", diff)
	}
}
