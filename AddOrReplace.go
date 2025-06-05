package cache

import "time"

// Adds a value to the cache or replaces the existing value in the cache under the specified key.
//
// # Parameters
//
//	key TKey
//
// The key under which the value is stored in the cache.
//
//	value TValue
//
// The value added to the cache or to be replaced in the cache under the specified key.
//
//	expiresAfter time.Duration [OPTIONAL]
//
// Time after the value is removed from the cache.
// If the parameter is ommited the defaultExpieresAfter value is used. See cache.New method.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	err error
//
// ErrInvalidExpireAfter - when passed expiresAter is not greater than zero or is not equal to the special value - NeverExpire.
func (cache *Cache[TKey, TValue]) AddOrReplace(key TKey, value TValue, expiresAfter ...time.Duration) (err error) {
	if len(expiresAfter) > 0 && !isExpirationValid(expiresAfter[0]) {
		return ErrInvalidExpireAfter
	}
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if _, found := cache.items[key]; !found {
		return cache.add(key, value, expiresAfter...)
	}
	return cache.replace(key, value, expiresAfter...)
}
