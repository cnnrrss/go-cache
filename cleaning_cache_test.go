package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSelfCleanup(t *testing.T) {
	c := NewSelfCleanup(150 * time.Millisecond)
	c.Set("k", "v")
	_, found := c.Get("k")
	assert.True(t, found)
	c.RLock()
	assert.Equal(t, 1, len(c.store))
	c.RUnlock()
	time.Sleep(360 * time.Millisecond) // time needs to be more than double the expiration to ensure the cleaner ran
	_, found = c.Get("k")
	assert.False(t, found)

	c.RLock()
	assert.Equal(t, 0, len(c.store))
	c.RUnlock()
}
