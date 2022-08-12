package concurrency

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mString string

func (s mString) ToBytes() []byte {
	return []byte(s)
}

func TestShard(t *testing.T) {
	shardMap := NewShardedMap[mString, int](10)

	shardMap.Set("a", 1)
	shardMap.Set("b", 2)
	shardMap.Set("c", 3)

	result := []int{}
	for _, k := range shardMap.Keys() {
		v := shardMap.Get(k)
		result = append(result, v)
	}

	sort.Ints(result)
	values := []int{1, 2, 3}
	if diff := cmp.Diff(values, result); diff != "" {
		t.Errorf("Sharded mismatch: {-want; +got}: %s", diff)
	}
}
