package cache

import (
	"sync"
)

func NewCacheInmemSharded[V any](shardingType ShardingType) Cache[string, V] {
	info := getShardingInfo(shardingType)

	shards := []*cacheInmem[string, V]{}
	for i := 0; i < info.Shards; i++ {
		shards = append(
			shards,
			&cacheInmem[string, V]{
				m:    &sync.RWMutex{},
				data: map[string]*cacheInmemItem[V]{},
			},
		)
	}

	return &cacheInmemSharded[V]{
		info:   info,
		shards: shards,
	}
}

type cacheInmemSharded[V any] struct {
	info   ShardingInfo
	shards []*cacheInmem[string, V]
}

func (c *cacheInmemSharded[V]) Get(key string) (V, bool) {
	return c.getShard(key).Get(key)
}

func (c *cacheInmemSharded[V]) GetEx(key string, exp Expiration) (V, bool) {
	return c.getShard(key).GetEx(key, exp)
}

func (c *cacheInmemSharded[V]) Set(key string, value V, exp Expiration) {
	c.getShard(key).Set(key, value, exp)
}

func (c *cacheInmemSharded[V]) Del(key string) {
	c.getShard(key).Del(key)
}

func (c *cacheInmemSharded[V]) CheckExpired() {
	for _, shard := range c.shards {
		shard.CheckExpired()
	}
}

func (c *cacheInmemSharded[V]) getShard(key string) *cacheInmem[string, V] {
	index := c.info.HashFunc(key)
	return c.shards[index]
}
