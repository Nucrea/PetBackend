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
