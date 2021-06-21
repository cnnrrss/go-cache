package cache

import (
	"sync"
	"time"
)

type (
	// Cache is an in-memory key:value store where cached items can expire after
	// the expiration period upon lookup.
	Cache struct {
		sync.RWMutex
		expiration time.Duration
		store      map[string]cached
	}

	cached struct {
		expires time.Time
		value   interface{}
	}
)

// New returns a new Cache with a given value expiration duration.
func New(expiration time.Duration) *Cache {
	return &Cache{
		expiration: expiration,
		store:      map[string]cached{},
	}
}

// Get an item from the cache. Returns the value or nil, and a bool indicating
// whether the key was found.
func (c *Cache) Get(k string) (interface{}, bool) {
	c.RLock()
	cached, found := c.store[k]
	c.RUnlock()
	if !found || cached.expired() {
		return nil, false
	}
	return cached.value, found
}

// Set adds an item to the cache or replaces any existing item.
func (c *Cache) Set(k string, v interface{}) {
	t := time.Now().Add(c.expiration)
	c.Lock()
	c.store[k] = cached{
		expires: t,
		value:   v,
	}
	c.Unlock()
}

func (c cached) expired() bool {
	return c.expires.Before(time.Now())
}
