package cache

import "github.com/thereisnoplanb/cache/enums/RemovalSource"

// Removes all values from the cache.
//
// # Remarks
//
// RemoveAll raises ItemRemoved with RemovalSource.Manual once for each removed value.
// The order of the events is unspecified.
func (cache *Cache[TKey, TValue]) RemoveAll() {
	cache.mutex.Lock()
	itemsRemovedEventArgs := make([]ItemRemovedEventArgs[TKey, TValue], 0, len(cache.items))
	defer func() {
		cache.mutex.Unlock()
		for _, itemRemovedEventArgs := range itemsRemovedEventArgs {
			cache.onItemRemoved(cache, itemRemovedEventArgs)
		}
	}()
	for key, item := range cache.items {
		if item.Timer != nil {
			item.Timer.Stop()
		}
		delete(cache.items, key)
		itemsRemovedEventArgs = append(itemsRemovedEventArgs, ItemRemovedEventArgs[TKey, TValue]{
			Key:           key,
			Value:         item.Value,
			RemovalSource: RemovalSource.Manual,
		})
	}
}
