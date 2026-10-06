package cache

import "time"

// Tries to reset the expiration duration in the cache item under the specified key.
//
// # Parameters
//
//	key TKey
//
// The key under which the value is stored in the cache.
//
// expiresAfter time.Duration [OPTIONAL]
//
// The duration after which the value is removed from the cache.
// If omitted, the default expiration specified in New or Must is used.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	success bool
//
// True if the reset operation was successful; false otherwise.
func (cache *Cache[TKey, TValue]) TryReset(key TKey, expiresAfter ...time.Duration) (success bool) {
	return cache.Reset(key, expiresAfter...) == nil
}
