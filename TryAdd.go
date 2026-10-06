package cache

import "time"

// Tries to add a value to the cache under the specified key.
//
// # Parameters
//
//	key TKey
//
// The key under which the value is stored in the cache.
//
//	value TValue
//
// The value added to the cache under the specified key.
//
//	expiresAfter time.Duration [OPTIONAL]
//
// The duration after which the value is removed from the cache.
// If omitted, the default expiration specified in New or Must is used.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	success bool
//
// True if the value was added; false otherwise.
func (cache *Cache[TKey, TValue]) TryAdd(key TKey, value TValue, expiresAfter ...time.Duration) (success bool) {
	return cache.Add(key, value, expiresAfter...) == nil
}
