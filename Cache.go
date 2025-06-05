package cache

import (
	"sync"
	"time"
)

// Indicates that a cached item should not be automatically removed.
const NeverExpire time.Duration = -1

// Cache.
type Cache[TKey comparable, TValue any] struct {
	items              map[TKey]item[TValue]
	mutex              *sync.Mutex
	defaultExpireAfter time.Duration
}

type item[TValue any] struct {
	Value TValue
	Timer *time.Timer
}
