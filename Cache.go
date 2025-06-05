package cache

import (
	"sync"
	"time"

	"github.com/thereisnoplanb/event"
)

// Indicates that a cached item should not be automatically removed.
const NeverExpire time.Duration = -1

// Cache.
type Cache[TKey comparable, TValue any] struct {
	items                map[TKey]item[TValue]
	mutex                *sync.Mutex
	defaultExpieresAfter time.Duration
	// Occurs when a value is added to the cache.
	ItemAdded event.Event[Cache[TKey, TValue], AddedEventArgs[TKey, TValue]]
	// Occurs when a value is replaced in the cache.
	ItemReplaced event.Event[Cache[TKey, TValue], ReplacedEventArgs[TKey, TValue]]
	// Occurs when a value is expired and removed from the cache.
	ItemExpired event.Event[Cache[TKey, TValue], ExpiredEventArgs[TKey, TValue]]
	// Occurs when a value is removed from the cache.
	ItemRemoved event.Event[Cache[TKey, TValue], RemovedEventArgs[TKey, TValue]]
}

type item[TValue any] struct {
	Value TValue
	Timer *time.Timer
}

// Provides data for Cache ItemAdded event.
type AddedEventArgs[TKey comparable, TValue any] struct {
	// The key under which the value is stored in the cache.
	Key TKey
	// The value added to the cache under the specified key.
	Value TValue
}

// Provides data for Cache ItemReplaced event.
type ReplacedEventArgs[TKey comparable, TValue any] struct {
	// The key under which the value is stored in the cache.
	Key TKey
	// The old value replaced in the cache by the new value at the specified key.
	OldValue TValue
	// The new value replacing the old value at the cache with the specified key.
	NewValue TValue
}

// Provides data for Cache ItemExpired event.
type ExpiredEventArgs[TKey comparable, TValue any] struct {
	// The key under which the value was stored in the cache.
	Key TKey
	// The value that has expierd in cache.
	Value TValue
}

// Provides data for Cache ItemRemoved event.
type RemovedEventArgs[TKey comparable, TValue any] struct {
	// The key under which the value was stored in the cache.
	Key TKey
	// The value that has been removed from cache.
	Value TValue
}
