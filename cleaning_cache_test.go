package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWithSelfCleanup(t *testing.T) {
	c := NewWithSelfCleanup(100 * time.Millisecond)
	c.Set("k", "v")

	_, found := c.Get("k")
	assert.True(t, found)

	c.RLock()
	assert.Equal(t, 1, len(c.store))
	c.RUnlock()

	time.Sleep(300 * time.Millisecond) // the time needs to be more than double the expiration to ensure the cleaner ran

	_, found = c.Get("k")
	assert.False(t, found)

	c.RLock()
	assert.Equal(t, 0, len(c.store))
	c.RUnlock()
}
