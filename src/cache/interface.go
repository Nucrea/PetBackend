package cache

import "time"

type Expiration struct {
	Ttl       time.Duration
	ExpiresAt time.Time
}

func (e Expiration) Get() time.Time {
	if e.Ttl != 0 {
		return time.Now().Add(e.Ttl)
	}
	return e.ExpiresAt
}

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	GetEx(key K, exp Expiration) (V, bool)

	Set(key K, value V, exp Expiration)

	Del(key K)
	CheckExpired(batchSize int)
}
