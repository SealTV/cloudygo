package concurrency

import (
	"sort"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFunnel(t *testing.T) {
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)
	channels := []chan int{c1, c2, c3}

	got := []int{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range Funnel(c1, c2, c3) {
			got = append(got, v)
		}
	}()

	for i, c := range channels {
		c <- i
	}
	for _, c := range channels {
		close(c)
	}

	wg.Wait()
	sort.Ints(got)
	if diff := cmp.Diff([]int{0, 1, 2}, got); diff != "" {
		t.Errorf("Funnel mismatch: {-want; +got}: %s", diff)
	}
}
