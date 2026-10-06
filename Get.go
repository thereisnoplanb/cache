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
// The value stored under the specified key, or the zero value of TValue if the key is absent.
//
//	found bool
//
// True if the specified key is present in the cache; false otherwise.
func (cache *Cache[TKey, TValue]) Get(key TKey) (value TValue, found bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if item, found := cache.items[key]; found {
		return item.Value, true
	}
	return value, false
}
