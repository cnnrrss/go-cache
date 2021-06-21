package cache

import (
	"sync"
	"time"
)

// Cache is an in-memory key-value store where keys expire after the configured period of time.
type Cache struct {
	sync.RWMutex
	expiration time.Duration
	store      map[string]cached
}

type cached struct {
	expires time.Time
	value   interface{}
}

// New returns a new Cache with a given value expiration duration.
func New(expiration time.Duration) *Cache {
	return &Cache{
		expiration: expiration,
		store:      map[string]cached{},
	}
}

// NewWithSelfCleanup returns a new Cache with a given expiration duration that
// automatically cleans it's internal storage once they are expired.
func NewWithSelfCleanup(expiration time.Duration) *Cache {
	c := New(expiration)

	go func() {
		for {
			select {
			case <-time.After(expiration * 2):

				c.RLock()
				copied := make(map[string]cached, len(c.store))
				for k, v := range c.store {
					copied[k] = v
				}
				c.RUnlock()

				var length int64
				for k, v := range copied {
					if v.expired() {
						c.Lock()
						delete(c.store, k)
						c.Unlock()
					} else {
						length++
					}
				}
			}
		}
	}()

	return c
}

// Get retrieves a value from cache or nil, and a bool indicating whether the key was found.
func (c *Cache) Get(k string) (interface{}, bool) {
	c.RLock()
	cached, found := c.store[k]
	c.RUnlock()
	if !found || cached.expired() {
		return nil, false
	}
	return cached.value, found
}

// Set adds a value to the cache or replaces an existing value.
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
