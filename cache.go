package cache

import "time"

type CacheItem struct {
	Value  string
	Expire time.Time
}

func (i CacheItem) isExpired() bool {
	if i.Expire.IsZero() {
		return false
	}

	if time.Now().After(i.Expire) {
		return true
	}

	return false
}

type Cache struct {
	storage map[string]CacheItem
}

func NewCache() Cache {
	return Cache{storage: make(map[string]CacheItem)}
}

func (c *Cache) Get(key string) (string, bool) {
	i, ok := c.storage[key]

	if !ok {
		return "", false
	}

	if i.isExpired() {
		delete(c.storage, key)

		return "", false
	}

	return i.Value, true
}

func (c *Cache) Put(key, value string) {
	c.storage[key] = CacheItem{
		Value:  value,
		Expire: time.Time{},
	}
}

func (c *Cache) Keys() []string {
	var keys []string

	for k, i := range c.storage {
		if !i.isExpired() {
			keys = append(keys, k)
		} else {
			delete(c.storage, k)
		}
	}
	
	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.storage[key] = CacheItem{
		Value:  value,
		Expire: deadline,
	}
}
