package cache

import "time"

// Creates a new cache with the specified key type and value type. Panics on error.
//
// # Parameters
//
//	expireAfter time.Duration
//
// The default value of time after the expired value is removed from the cache.
// The special value NeverExpire can be used to indicate that the values are stored in the cache until manually removed.
//
// # Returns
//
//	cache *Cache[TKey, TValue]
//
// A cache with the specified key type and value type.
//
//	# Remarks
//
// Panics when expireAfter time.Duration param is not grater than zero or is not equal to the special value - NeverExpire.
func Must[TKey comparable, TValue any](expireAfter time.Duration) (cache *Cache[TKey, TValue]) {
	c, err := New[TKey, TValue](expireAfter)
	if err != nil {
		panic(err)
	}
	return c
}
