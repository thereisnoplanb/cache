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
// Time after the value is removed from the cache.
// If the parameter is ommited the defaultExpieresAfter value is used. See cache.New method.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	success bool
//
// True when adding operation is succesful; false otherwise.
func (cache *Cache[TKey, TValue]) TryAdd(key TKey, value TValue, expiresAfter ...time.Duration) (success bool) {
	return cache.Add(key, value, expiresAfter...) == nil
}
