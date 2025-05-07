// The sharded map is inspired by this article: https://medium.com/@isurucuma/let-us-build-a-thread-safe-shardedmap-in-golang-273b0c6c092b
package memheapv2

import (
	"hash/fnv"
	"sync"
)

type shard[T any] struct {
	sync.Mutex
	m map[string]T
}

type shardedMap[T any] []*shard[T]

func newShardedMap[T any](nShards uint) shardedMap[T] {
	shards := make([]*shard[T], nShards)

	for i := range nShards {
		shardMap := make(map[string]T)
		shards[i] = &shard[T]{m: shardMap}
	}

	return shards
}

func (sm shardedMap[T]) getShardIndex(key string) int {
	if len(sm) == 1 {
		return 0
	}
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % len(sm)
}

func (sm shardedMap[T]) getShard(key string) *shard[T] {
	idx := sm.getShardIndex(key)
	return sm[idx]
}

func (sm shardedMap[T]) get(key string) T {
	shard := sm.getShard(key)
	return shard.m[key]
}

func (sm shardedMap[T]) itemCount() int {
	total := 0
	for _, shard := range sm {
		total += len(shard.m)
	}
	return total
}
