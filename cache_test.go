package cache

import (
	"testing"
	"time"
)

var fixture = New(5 * time.Minute)

func TestSet_add(t *testing.T) {
	fixture.Set("k", 1)
	if len(fixture.store) != 1 {
		t.Error("expected 1 entry in the cache after add")
	}
}

func TestSet_replace(t *testing.T) {
	fixture.Set("k", 2)
	if len(fixture.store) != 1 {
		t.Error("expected 1 entry in the cache after replace")
	}
}

func TestGet_found(t *testing.T) {
	v, found := fixture.Get("k")
	if !found || v.(int) != 2 {
		t.Error("expected to find existing entry")
	}
}

func TestGet_missing(t *testing.T) {
	_, found := fixture.Get("j")
	if found {
		t.Error("expected to not find value")
	}
}

func TestGet_expired(t *testing.T) {
	c := New(1 * time.Nanosecond)
	c.Set("k", 1)
	time.Sleep(2 * time.Nanosecond)
	_, found := c.Get("k")
	if found {
		t.Error("expected to not find expired value")
	}
}
