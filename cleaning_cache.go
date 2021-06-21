package cache

import (
	"expvar"
	"fmt"
	"time"
)

// NewWithSelfCleanup returns a new Cache with a given expiration duration that
// automatically cleans it's internal storage once they are expired.
func NewWithSelfCleanup(expiration time.Duration) *Cache {
	c := New(expiration)
	before := expvar.NewInt(fmt.Sprintf("size_before_%p", c))
	after := expvar.NewInt(fmt.Sprintf("size_after_%p", c))
	duration := expvar.NewInt(fmt.Sprintf("clean_duration_nanos_%p", c))
	go func() {
		for {
			select {
			case <-time.After(expiration * 2):
				start := time.Now()

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

				before.Set(int64(len(copied)))
				after.Set(length)
				duration.Set(time.Now().Sub(start).Nanoseconds())
			}
		}
	}()
	return c
}
