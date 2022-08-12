package concurrency

import (
	"crypto/sha1"
	"sync"
)

type Key interface {
	ToBytes() []byte
	comparable
}

type Shard[K Key, V any] struct {
	sync.RWMutex
	m map[K]V
}

type ShardedMap[K Key, V any] []*Shard[K, V]

func NewShardedMap[K Key, V any](nshards int) ShardedMap[K, V] {
	shards := make([]*Shard[K, V], 0, nshards)
	for i := 0; i < nshards; i++ {
		shards = append(shards, &Shard[K, V]{m: map[K]V{}})
	}

	return shards
}

func (m ShardedMap[K, V]) getShardIndex(key K) int {
	checksum := sha1.Sum(key.ToBytes())
	hash := int(checksum[2])<<56 |
		int(checksum[4])<<48 |
		int(checksum[6])<<40 |
		int(checksum[8])<<32 |
		int(checksum[10])<<24 |
		int(checksum[12])<<16 |
		int(checksum[14])<<8 |
		int(checksum[16])
	if hash < 0 {
		hash = -hash
	}
	return hash % len(m)
}

func (m ShardedMap[K, V]) getShard(key K) *Shard[K, V] {
	index := m.getShardIndex(key)
	return m[index]
}

func (m ShardedMap[K, V]) Get(key K) V {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (m ShardedMap[K, V]) Set(key K, val V) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = val
}

func (m ShardedMap[K, V]) Keys() []K {
	keys := make([]K, 0)

	mtx := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(m))

	for _, shard := range m {
		go func(s *Shard[K, V]) {
			defer wg.Done()

			s.RLock()
			defer s.RUnlock()

			for key := range s.m {
				mtx.Lock()
				keys = append(keys, key)
				mtx.Unlock()
			}
		}(shard)
	}

	wg.Wait()

	return keys
}
