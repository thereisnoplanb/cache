package cache

import (
	"sync"
	"time"

	"github.com/thereisnoplanb/event"
)

// Indicates that a cached item should not be automatically removed.
const NeverExpire time.Duration = -1

// Cache represents a generic in-memory cache with support for item expiration and events.
type Cache[TKey comparable, TValue any] struct {
	items              map[TKey]item[TValue]
	mutex              *sync.Mutex
	defaultExpireAfter time.Duration
	ItemAdded          event.Event[Cache[TKey, TValue], ItemAddedEventArgs[TKey, TValue]]
	onItemAdded        event.Invoke[Cache[TKey, TValue], ItemAddedEventArgs[TKey, TValue]]
	ItemRemoved        event.Event[Cache[TKey, TValue], ItemRemovedEventArgs[TKey, TValue]]
	onItemRemoved      event.Invoke[Cache[TKey, TValue], ItemRemovedEventArgs[TKey, TValue]]
	ItemReplaced       event.Event[Cache[TKey, TValue], ItemReplacedEventArgs[TKey, TValue]]
	onItemReplaced     event.Invoke[Cache[TKey, TValue], ItemReplacedEventArgs[TKey, TValue]]
}

type item[TValue any] struct {
	Value           TValue
	Timer           *time.Timer
	expirationToken *expirationToken
}

type expirationToken struct {
	marker byte
}
