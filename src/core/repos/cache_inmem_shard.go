package repos

import (
	"sync"
)

type ShardingType int

const (
	ShardingTypeJWT ShardingType = iota
	ShardingTypeInteger
)

type shardingHashFunc func(key string) int

func getShardingInfo(shardingType ShardingType) (int, shardingHashFunc) {
	switch shardingType {
	case ShardingTypeInteger:
		return 10, func(key string) int {
			char := int(key[len(key)-1])
			return char - 0x30
		}
	case ShardingTypeJWT:
		return 36, func(key string) int {
			char := int(key[len(key)-1])
			if char >= 0x30 && char <= 0x39 {
				return char - 0x30
			}
			if char >= 0x41 && char <= 0x5A {
				return char - 0x41
			}
			return char - 0x61
		}
	}

	return 1, func(key string) int {
		return 0
	}
}

func NewCacheInmemSharded[V any](defaultTtlSeconds int, shardingType ShardingType) Cache[string, V] {
	shards, hashFunc := getShardingInfo(shardingType)

	inmems := []*cacheInmem[string, V]{}
	for i := 0; i < shards; i++ {
		inmems = append(
			inmems,
			&cacheInmem[string, V]{
				m:          &sync.Mutex{},
				data:       map[string]*cacheInmemItem[V]{},
				ttlSeconds: defaultTtlSeconds,
			},
		)
	}

	return &cacheInmemSharded[V]{
		shards:   inmems,
		hashFunc: hashFunc,
	}
}

type cacheInmemSharded[V any] struct {
	hashFunc shardingHashFunc
	shards   []*cacheInmem[string, V]
}

func (c *cacheInmemSharded[V]) Get(key string) (V, bool) {
	return c.getShard(key).Get(key)
}

func (c *cacheInmemSharded[V]) Set(key string, value V, ttlSeconds int) {
	c.getShard(key).Set(key, value, ttlSeconds)
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
	index := c.hashFunc(key)
	return c.shards[index]
}
