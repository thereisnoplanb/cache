package cache

// Gets a value from the cache under the specified key.
//
// # Parameters
//
//	key TKey
//
// The key under which the value is stored in the cache.
//
// # Returns
//
//	value TValue
//
// Value from the specified key.
//
//	found bool
//
// True if specified key is present in cache; false otherwise.
func (cache *Cache[TKey, TValue]) Get(key TKey) (value TValue, found bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if item, found := cache.items[key]; found {
		return item.Value, true
	}
	return value, false
}
