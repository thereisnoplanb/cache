package cache

import "time"

// Creates a new cache with the specified default expiration duration. Panics if the duration is invalid.
//
// # Parameters
//
//	expireAfter time.Duration
//
// The default duration after which each value is removed from the cache.
// The special value NeverExpire disables automatic expiration by default.
//
// # Returns
//
//	cache *Cache[TKey, TValue]
//
// A cache with the specified key type and value type.
//
// # Remarks
//
// Panics when expireAfter is neither greater than zero nor equal to NeverExpire.
func Must[TKey comparable, TValue any](expireAfter time.Duration) (cache *Cache[TKey, TValue]) {
	c, err := New[TKey, TValue](expireAfter)
	if err != nil {
		panic(err)
	}
	return c
}
