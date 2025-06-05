package cache

// Removes all values from the cache.
func (cache *Cache[TKey, TValue]) RemoveAll() {
	cache.mutex.Lock()
	items := cache.items
	defer func() {
		cache.mutex.Unlock()
		for key, item := range items {
			cache.ItemRemoved.Invoke(cache, RemovedEventArgs[TKey, TValue]{
				Key:   key,
				Value: item.Value,
			})
		}
	}()
	for key, item := range cache.items {
		if item.Timer != nil && !item.Timer.Stop() {
			<-item.Timer.C
		}
		delete(cache.items, key)
	}
}
