package cache

// Removes all values from the cache.
func (cache *Cache[TKey, TValue]) RemoveAll() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	for key, item := range cache.items {
		if item.Timer != nil && !item.Timer.Stop() {
			<-item.Timer.C
		}
		delete(cache.items, key)
	}
}
