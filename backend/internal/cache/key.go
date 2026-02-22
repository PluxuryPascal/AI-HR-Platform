package cache

import (
	"time"
)

type Key[T any] struct {
	namespace string
	ttl       time.Duration
}

func NewKey[T any](namespace string) Key[T] {
	return Key[T]{
		namespace: namespace,
	}
}

func (k Key[T]) WithTTL(ttl time.Duration) Key[T] {
	k.ttl = ttl
	return k
}

// -----------------------------------------------------
// Реестр ключей кэша
// -----------------------------------------------------

var (
	SessionKey   = NewKey[string]("session")
	RateLimitKey = NewKey[int64]("rate_limit")
)
