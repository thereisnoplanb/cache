package cache

import "time"

// Replaces the value in the cache under the specified key.
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
// If the parameter is ommited, the defaultExpieresAfter value is used. See cache.New method.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	err error
//
// ErrKeyNotFound - when the specified key does not exist in the cache.
// ErrInvalidExpireAfter - when passed expiresAter is not greater than zero or not equal to the special value - NeverExpire.
func (cache *Cache[TKey, TValue]) Replace(key TKey, value TValue, expiresAfter ...time.Duration) (err error) {
	if len(expiresAfter) > 0 && !isExpirationValid(expiresAfter[0]) {
		return ErrInvalidExpireAfter
	}
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if _, found := cache.items[key]; !found {
		return ErrKeyNotFound
	}
	return cache.replace(key, value, expiresAfter...)
}

func (cache *Cache[TKey, TValue]) replace(key TKey, value TValue, expiresAfter ...time.Duration) (err error) {
	if item, found := cache.items[key]; found {
		if item.Timer != nil && !item.Timer.Stop() {
			<-item.Timer.C
		}
	}
	expiration := cache.defaultExpieresAfter
	if len(expiresAfter) > 0 && (expiresAfter[0] > 0 || expiresAfter[0] == NeverExpire) {
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
