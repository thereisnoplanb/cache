package cache

import (
	"runtime"
	"sync"
	"time"

	"github.com/thereisnoplanb/cache/enums/RemovalSource"
	"github.com/thereisnoplanb/event"
)

type ItemAddedEventArgs[TKey comparable, TValue any] struct {
	Key   TKey
	Value TValue
}

type ItemRemovedEventArgs[TKey comparable, TValue any] struct {
	Key           TKey
	Value         TValue
	RemovalSource RemovalSource.Enum
}

type ItemReplacedEventArgs[TKey comparable, TValue any] struct {
	Key      TKey
	OldValue TValue
	NewValue TValue
}

// Creates a new cache with the specified default expiration duration.
//
// # Parameters
//
//	expireAfter time.Duration
//
// The default duration after which each value is removed from the cache.
// The special value NeverExpire disables automatic expiration by default.
//
// # Returns
//
//	cache *Cache[TKey, TValue]
//
// A cache with the specified key type and value type.
//
//	err error
//
// Errors:
//   - ErrInvalidExpireAfter - when expireAfter is neither greater than zero nor equal to NeverExpire.
func New[TKey comparable, TValue any](expireAfter time.Duration) (c *Cache[TKey, TValue], err error) {
	if !isExpirationValid(expireAfter) {
		return nil, ErrInvalidExpireAfter
	}
	itemAdded, onItemAdded := event.New[Cache[TKey, TValue], ItemAddedEventArgs[TKey, TValue]]()
	itemRemoved, onItemRemoved := event.New[Cache[TKey, TValue], ItemRemovedEventArgs[TKey, TValue]]()
	itemReplaced, onItemReplaced := event.New[Cache[TKey, TValue], ItemReplacedEventArgs[TKey, TValue]]()
	c = &Cache[TKey, TValue]{
		items:              make(map[TKey]item[TValue]),
		mutex:              &sync.Mutex{},
		defaultExpireAfter: expireAfter,
		ItemAdded:          itemAdded,
		onItemAdded:        onItemAdded,
		ItemRemoved:        itemRemoved,
		onItemRemoved:      onItemRemoved,
		ItemReplaced:       itemReplaced,
		onItemReplaced:     onItemReplaced,
	}
	runtime.SetFinalizer(c, finalize[TKey, TValue])
	return c, nil
}

func finalize[TKey comparable, TValue any](cache *Cache[TKey, TValue]) {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.Timer != nil {
			item.Timer.Stop()
		}
		delete(cache.items, key)
	}
	cache.mutex.Unlock()
}

func isExpirationValid(expiration time.Duration) bool {
	return expiration > 0 || expiration == NeverExpire
}
