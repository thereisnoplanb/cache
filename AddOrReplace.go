package cache

import "time"

// Adds a value to the cache or replaces the value in the cache under the specified key.
//
// # Parameters
//
//	key TKey
//
// The key under which the value is stored in the cache.
//
//	value TValue
//
// The value added to the cache or replaced in the cache under the specified key.
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
//   - ErrInvalidExpireAfter - when the specified expiration is neither greater than zero nor equal to NeverExpire.
//
// # Remarks
//
// On success, AddOrReplace raises ItemAdded if the key was absent, or ItemReplaced if it was present.
func (cache *Cache[TKey, TValue]) AddOrReplace(key TKey, value TValue, expiresAfter ...time.Duration) (err error) {
	if len(expiresAfter) > 0 && !isExpirationValid(expiresAfter[0]) {
		return ErrInvalidExpireAfter
	}
	cache.mutex.Lock()
	added := false
	replaced := false
	var oldValue TValue
	defer func() {
		cache.mutex.Unlock()
		if added {
			cache.onItemAdded(cache, ItemAddedEventArgs[TKey, TValue]{
				Key:   key,
				Value: value,
			})
		}
		if replaced {
			cache.onItemReplaced(cache, ItemReplacedEventArgs[TKey, TValue]{
				Key:      key,
				OldValue: oldValue,
				NewValue: value,
			})
		}
	}()
	if item, found := cache.items[key]; !found {
		cache.add(key, value, expiresAfter...)
		added = true
		return nil
	} else {
		oldValue = item.Value
	}
	cache.replace(key, value, expiresAfter...)
	replaced = true
	return nil
}
