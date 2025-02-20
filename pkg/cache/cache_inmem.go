package cache

import (
	"sync"
	"time"
)

func NewCacheInmem[K comparable, V any]() Cache[K, V] {
	return &cacheInmem[K, V]{
		m:    &sync.RWMutex{},
		data: map[K]*cacheInmemItem[V]{},
	}
}

type cacheInmemItem[T any] struct {
	Value      T
	Expiration time.Time
}

type cacheInmem[K comparable, V any] struct {
	m    *sync.RWMutex
	data map[K]*cacheInmemItem[V]
}

func (c *cacheInmem[K, V]) Get(key K) (V, bool) {
	c.m.RLock()
	defer c.m.RUnlock()

	var v V

	item, ok := c.data[key]
	if !ok {
		return v, false
	}
	if time.Now().Before(item.Expiration) {
		return item.Value, true
	}

	return v, false
}

func (c *cacheInmem[K, V]) GetEx(key K, exp Expiration) (V, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	item, ok := c.data[key]
	if !ok {
		var v V
		return v, false
	}

	if time.Now().Before(item.Expiration) {
		c.data[key].Expiration = exp.Get()
		return item.Value, true
	}

	delete(c.data, key)

	var v V
	return v, false
}

func (c *cacheInmem[K, V]) Set(key K, value V, exp Expiration) {
	c.m.Lock()
	defer c.m.Unlock()

	item := &cacheInmemItem[V]{
		Value:      value,
		Expiration: exp.Get(),
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

	for key, item := range c.data {
		if time.Now().After(item.Expiration) {
			delete(c.data, key)
		}
	}
}
