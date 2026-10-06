package cache

import (
	"time"
)

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
//	expiresAfter ...time.Duration [OPTIONAL]
//
// The duration after which the value is removed from the cache.
// If omitted, the default expiration specified in New or Must is used.
// If multiple values are provided, only the first is used.
// The special value NeverExpire can be used to indicate that the value is cached until manually removed.
//
// # Returns
//
//	err error
//
// Errors:
//   - ErrKeyNotFound - when the specified key does not exist in the cache.
//   - ErrInvalidExpireAfter - when the specified expiration is neither greater than zero nor equal to NeverExpire.
//
// # Remarks
//
// On success, Replace raises the ItemReplaced event after updating the value.
func (cache *Cache[TKey, TValue]) Replace(key TKey, value TValue, expiresAfter ...time.Duration) (err error) {
	if len(expiresAfter) > 0 && !isExpirationValid(expiresAfter[0]) {
		return ErrInvalidExpireAfter
	}
	cache.mutex.Lock()
	var oldValue item[TValue]
	found := false
	replaced := false
	defer func() {
		cache.mutex.Unlock()
		if replaced {
			cache.onItemReplaced(cache, ItemReplacedEventArgs[TKey, TValue]{
				Key:      key,
				OldValue: oldValue.Value,
				NewValue: value,
			})
		}
	}()
	oldValue, found = cache.items[key]
	if !found {
		return ErrKeyNotFound
	}
	cache.replace(key, value, expiresAfter...)
	replaced = true
	return nil
}

func (cache *Cache[TKey, TValue]) replace(key TKey, value TValue, expiresAfter ...time.Duration) {
	if item, found := cache.items[key]; found {
		if item.Timer != nil {
			item.Timer.Stop()
		}
	}
	expiration := cache.defaultExpireAfter
	if len(expiresAfter) > 0 && (expiresAfter[0] > 0 || expiresAfter[0] == NeverExpire) {
		expiration = expiresAfter[0]
	}
	timer, token := cache.scheduleExpiration(key, expiration)
	cache.items[key] = item[TValue]{
		Value:           value,
		Timer:           timer,
		expirationToken: token,
	}
}
