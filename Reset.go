package cache

import "time"

// Resets the expiration duration for the cache item under the specified key.
//
// # Parameters
//
// key TKey
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
// err error
//
// Errors:
//   - ErrKeyNotFound - when the specified key does not exist in the cache.
//   - ErrInvalidExpireAfter - when the specified expiration is neither greater than zero nor equal to NeverExpire.
//
// # Remarks
//
// Reset does not raise an event immediately. ItemRemoved is raised if the item later expires.
func (cache *Cache[TKey, TValue]) Reset(key TKey, expiresAfter ...time.Duration) (err error) {
	if len(expiresAfter) > 0 && !isExpirationValid(expiresAfter[0]) {
		return ErrInvalidExpireAfter
	}
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if item, found := cache.items[key]; found {
		if item.Timer != nil {
			item.Timer.Stop()
		}
		cache.reset(key, expiresAfter...)
		return nil
	}
	return ErrKeyNotFound
}

func (cache *Cache[TKey, TValue]) reset(key TKey, expiresAfter ...time.Duration) {
	expiration := cache.defaultExpireAfter
	if len(expiresAfter) > 0 {
		expiration = expiresAfter[0]
	}
	timer, token := cache.scheduleExpiration(key, expiration)
	item := cache.items[key]
	item.Timer = timer
	item.expirationToken = token
	cache.items[key] = item
}
