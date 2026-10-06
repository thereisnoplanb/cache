package cache

import (
	"time"

	"github.com/thereisnoplanb/cache/enums/RemovalSource"
)

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
// The duration after which the value is removed from the cache.
// If omitted, the default expiration specified in New or Must is used.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	err error
//
// Errors:
//   - ErrKeyAlreadyExists - when the specified key already exists in the cache.
//   - ErrInvalidExpireAfter - when the specified expiration is neither greater than zero nor equal to NeverExpire.
//
// # Remarks
//
// On success, Add raises the ItemAdded event after adding the value to the cache.
func (cache *Cache[TKey, TValue]) Add(key TKey, value TValue, expiresAfter ...time.Duration) (err error) {
	if len(expiresAfter) > 0 && !isExpirationValid(expiresAfter[0]) {
		return ErrInvalidExpireAfter
	}
	cache.mutex.Lock()
	added := false
	defer func() {
		cache.mutex.Unlock()
		if added {
			cache.onItemAdded(cache, ItemAddedEventArgs[TKey, TValue]{
				Key:   key,
				Value: value,
			})
		}
	}()
	if _, found := cache.items[key]; found {
		return ErrKeyAlreadyExists
	}
	cache.add(key, value, expiresAfter...)
	added = true
	return nil
}

func (cache *Cache[TKey, TValue]) add(key TKey, value TValue, expiresAfter ...time.Duration) {
	expiration := cache.defaultExpireAfter
	if len(expiresAfter) > 0 {
		expiration = expiresAfter[0]
	}
	timer, token := cache.scheduleExpiration(key, expiration)
	cache.items[key] = item[TValue]{
		Value:           value,
		Timer:           timer,
		expirationToken: token,
	}
}

func (cache *Cache[TKey, TValue]) scheduleExpiration(key TKey, expiration time.Duration) (*time.Timer, *expirationToken) {
	if expiration == NeverExpire {
		return nil, nil
	}
	token := &expirationToken{}
	timer := time.AfterFunc(expiration, func() {
		cache.expire(key, token)
	})
	return timer, token
}

func (cache *Cache[TKey, TValue]) expire(key TKey, token *expirationToken) {
	cache.mutex.Lock()
	removed := false
	var expiredItem item[TValue]
	defer func() {
		cache.mutex.Unlock()
		if removed {
			cache.onItemRemoved(cache, ItemRemovedEventArgs[TKey, TValue]{
				Key:           key,
				Value:         expiredItem.Value,
				RemovalSource: RemovalSource.Expiration,
			})
		}
	}()
	item, found := cache.items[key]
	if !found || item.expirationToken != token {
		return
	}
	expiredItem = item
	delete(cache.items, key)
	removed = true
}
