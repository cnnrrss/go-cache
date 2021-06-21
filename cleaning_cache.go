package cache

import (
	"time"
)

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
