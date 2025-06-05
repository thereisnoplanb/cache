package cache

import "time"

// Tries to replace the value in the cache under the specified key.
//
// # Parameters
//
//	key TKey
//
// The key under which the value is stored in the cache.
//
//	value TValue
//
// The value to be replaced in the cache under the specified key.
//
//	expiresAfter time.Duration [OPTIONAL]
//
// The time after the value is removed from the cache.
// If the parameter is ommited, the default expireAfter value is used. See cache.New, cache.Must methods.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	success bool
//
// True when replacing operation is succesful; false otherwise.
func (cache *Cache[TKey, TValue]) TryReplace(key TKey, value TValue, expiresAfter ...time.Duration) (success bool) {
	return cache.Replace(key, value, expiresAfter...) == nil
}
