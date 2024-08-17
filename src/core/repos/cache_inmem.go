package repos

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V, ttlSeconds int)
	Del(key K)
	CheckExpired()
}

func NewCacheInmem[K comparable, V any](ttlSeconds int) Cache[K, V] {
	return &cacheInmem[K, V]{
		m:          &sync.Mutex{},
		data:       map[K]*cacheInmemItem[V]{},
		ttlSeconds: ttlSeconds,
	}
}

type cacheInmemItem[T any] struct {
	Value      T
	Ttl        int64
	Expiration int64
}

type cacheInmem[K comparable, V any] struct {
	m          *sync.Mutex
	data       map[K]*cacheInmemItem[V]
	ttlSeconds int
}

func (c *cacheInmem[K, V]) Get(key K) (V, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	item, ok := c.data[key]
	if !ok {
		var v V
		return v, false
	}

	timestamp := time.Now().Unix()
	if item.Expiration > timestamp {
		item.Expiration = timestamp + item.Ttl
		return item.Value, true
	}

	delete(c.data, key)

	var v V
	return v, false
}

func (c *cacheInmem[K, V]) Set(key K, value V, ttlSeconds int) {
	c.m.Lock()
	defer c.m.Unlock()

	ttl := int64(c.ttlSeconds)

	expiration := time.Now().Unix() + ttl
	item := &cacheInmemItem[V]{
		Value:      value,
		Ttl:        ttl,
		Expiration: expiration,
	}
	c.data[key] = item
}

func (c *cacheInmem[K, V]) Del(key K) {
	c.m.Lock()
	defer c.m.Unlock()

	delete(c.data, key)
}

func (c *cacheInmem[K, V]) CheckExpired() {
	if len(c.data) == 0 {
		return
	}

	c.m.Lock()
	defer c.m.Unlock()

	itemsToProcess := 1000
	for key, item := range c.data {
		timestamp := time.Now().Unix()
		if item.Expiration <= timestamp {
			delete(c.data, key)
		}

		itemsToProcess--
		if itemsToProcess <= 0 {
			return
		}
	}
}
