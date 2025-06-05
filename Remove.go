package cache

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
// If passed key value is not present is cache no action is performed.
func (cache *Cache[TKey, TValue]) Remove(key TKey) {
	cache.mutex.Lock()
	item, found := cache.items[key]
	defer func() {
		cache.mutex.Unlock()
		if found {
			cache.ItemRemoved.Invoke(cache, RemovedEventArgs[TKey, TValue]{
				Key:   key,
				Value: item.Value,
			})
		}
	}()
	if found {
		if item.Timer != nil && !item.Timer.Stop() {
			<-item.Timer.C
		}
		delete(cache.items, key)
	}
}
