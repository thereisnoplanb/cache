package cache

import (
	"runtime"
	"sync"
	"time"
)

// Creates a new cache with the specified key type and value type.
//
// # Parameters
//
//	expireAfter time.Duration
//
// The default value of time after the expired value is removed from the cache.
// The special value NeverExpire can be used to indicate that the values are stored in the cache until manually removed.
//
// # Returns
//
//	cache *Cache[TKey, TValue]
//
// A cache with the specified key type and value type.
//
//	err error
//
// ErrInvalidExpireAfter - when expireAfter parameter is not grater than zero or not equal to the special value - NeverExpire.
func New[TKey comparable, TValue any](expireAfter time.Duration) (cache *Cache[TKey, TValue], err error) {
	if !isExpirationValid(expireAfter) {
		return nil, ErrInvalidExpireAfter
	}
	cache = &Cache[TKey, TValue]{
		items:              make(map[TKey]item[TValue]),
		mutex:              &sync.Mutex{},
		defaultExpireAfter: expireAfter,
	}
	runtime.SetFinalizer(cache, finalize[TKey, TValue])
	return cache, nil
}

func finalize[TKey comparable, TValue any](cache *Cache[TKey, TValue]) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	for key, item := range cache.items {
		if item.Timer != nil && !item.Timer.Stop() {
			<-item.Timer.C
		}
		delete(cache.items, key)
	}
}

func isExpirationValid(expiration time.Duration) bool {
	return expiration > 0 || expiration == NeverExpire
}
