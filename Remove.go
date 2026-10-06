package cache

import "github.com/thereisnoplanb/cache/enums/RemovalSource"

// Removes the value from the cache under the specified key.
//
// # Parameters
//
//	key TKey
//
// The key under which the value is stored in the cache.
//
// # Remarks
//
// If the specified key is not present in the cache, no action is performed.
// If it is present, Remove raises ItemRemoved with RemovalSource.Manual after removing the value.
func (cache *Cache[TKey, TValue]) Remove(key TKey) {
	cache.mutex.Lock()
	var itemRemovedEventArgs ItemRemovedEventArgs[TKey, TValue]
	removed := false
	defer func() {
		cache.mutex.Unlock()
		if removed {
			cache.onItemRemoved(cache, itemRemovedEventArgs)
		}
	}()
	if item, found := cache.items[key]; found {
		if item.Timer != nil {
			item.Timer.Stop()
		}
		delete(cache.items, key)
		itemRemovedEventArgs = ItemRemovedEventArgs[TKey, TValue]{
			Key:           key,
			Value:         item.Value,
			RemovalSource: RemovalSource.Manual,
		}
		removed = true
	}
}
