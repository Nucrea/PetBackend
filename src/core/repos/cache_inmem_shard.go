package repos

import (
	"sync"
)

func NewCacheInmemSharded[K comparable, V any](
	ttlSeconds, shards int,
	hashFunc func(key K) int,
) Cache[K, V] {
	inmems := []*cacheInmem[K, V]{}
	for i := 0; i < shards; i++ {
		inmems = append(
			inmems,
			&cacheInmem[K, V]{
				m:          &sync.Mutex{},
				data:       map[K]*cacheInmemItem[V]{},
				ttlSeconds: ttlSeconds,
			},
		)
	}

	return &cacheInmemSharded[K, V]{
		shards:   inmems,
		hashFunc: hashFunc,
	}
}

type cacheInmemSharded[K comparable, V any] struct {
	hashFunc func(key K) int
	shards   []*cacheInmem[K, V]
}

func (c *cacheInmemSharded[K, V]) Get(key K) (V, bool) {
	return c.getShard(key).Get(key)
}

func (c *cacheInmemSharded[K, V]) Set(key K, value V, ttlSeconds int) {
	c.getShard(key).Set(key, value, ttlSeconds)
}

func (c *cacheInmemSharded[K, V]) Del(key K) {
	c.getShard(key).Del(key)
}

func (c *cacheInmemSharded[K, V]) CheckExpired() {
	for _, shard := range c.shards {
		shard.CheckExpired()
	}
}

func (c *cacheInmemSharded[K, V]) getShard(key K) *cacheInmem[K, V] {
	index := c.hashFunc(key)
	return c.shards[index]
}
