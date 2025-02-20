package cache

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	GetEx(key K, exp Expiration) (V, bool)

	Set(key K, value V, exp Expiration)

	Del(key K)
	CheckExpired()
}
