package cache

import "time"

// Adds a value to the cache under the specified key.
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
// If the parameter is ommited the default expireAfter value is used. See cache.New, cache.Must methods.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	err error
//
// ErrKeyAlreadyExists - when the specified key already exists in the cache.
// ErrInvalidExpireAfter - when passed expiresAter is not greater than zero or is not equal to the special value - NeverExpire.
func (cache *Cache[TKey, TValue]) Add(key TKey, value TValue, expiresAfter ...time.Duration) (err error) {
	if len(expiresAfter) > 0 && !isExpirationValid(expiresAfter[0]) {
		return ErrInvalidExpireAfter
	}
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if _, found := cache.items[key]; found {
		return ErrKeyAlreadyExists
	}
	return cache.add(key, value, expiresAfter...)
}

func (cache *Cache[TKey, TValue]) add(key TKey, value TValue, expiresAfter ...time.Duration) (err error) {
	expiration := cache.defaultExpireAfter
	if len(expiresAfter) > 0 {
		expiration = expiresAfter[0]
	}
	var timer *time.Timer = nil
	if expiration != NeverExpire {
		timer = time.AfterFunc(expiration, func() {
			cache.mutex.Lock()
			defer cache.mutex.Unlock()
			delete(cache.items, key)
		})
	}
	cache.items[key] = item[TValue]{
		Value: value,
		Timer: timer,
	}
	return nil
}
