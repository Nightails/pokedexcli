package pokecache

import (
	"testing"
	"time"
)

// Test that Add inserts and Get retrieves the same value.
func TestAddAndGet(t *testing.T) {
	c := NewCache(time.Second)
	key := "cache:test:k1"
	val := []byte("bulbasaur")

	c.Add(key, val)

	got, ok := c.Get(key)
	if !ok {
		t.Fatalf("expected key to exist")
	}
	if string(got) != string(val) {
		t.Fatalf("got %q, want %q", string(got), string(val))
	}
}

// Test that Get returns false for missing keys.
func TestGetMissing(t *testing.T) {
	c := NewCache(time.Second)

	if _, ok := c.Get("cache:test:missing"); ok {
		t.Fatalf("expected ok=false for missing key")
	}
}

// Test that Add does not overwrite an existing entry with the same key.
func TestAddDoesNotOverwrite(t *testing.T) {
	c := NewCache(time.Second)
	key := "cache:test:k2"
	first := []byte("first")
	second := []byte("second")

	c.Add(key, first)
	c.Add(key, second) // should be ignored

	got, ok := c.Get(key)
	if !ok {
		t.Fatalf("expected key to exist")
	}
	if string(got) != string(first) {
		t.Fatalf("value was overwritten: got %q, want %q", string(got), string(first))
	}
}

// Test that reapLoop removes only entries older than the interval.
func TestReapLoop(t *testing.T) {
	interval := 50 * time.Millisecond
	c := NewCache(interval)

	oldKey := "cache:test:old"
	newKey := "cache:test:new"

	// Insert entries with controlled timestamps.
	c.Entries[oldKey] = cacheEntry{
		createdAt: time.Now().Add(-2 * interval),
		val:       []byte("old"),
	}
	c.Entries[newKey] = cacheEntry{
		createdAt: time.Now(),
		val:       []byte("new"),
	}

	c.reapLoop(interval)

	if _, ok := c.Get(oldKey); ok {
		t.Fatalf("expected old entry to be reaped")
	}
	if got, ok := c.Get(newKey); !ok || string(got) != "new" {
		t.Fatalf("expected new entry to remain; ok=%v got=%q", ok, string(got))
	}
}

// Ensure NewCache initializes an empty, usable cache.
func TestNewCacheInitializesCache(t *testing.T) {
	c := NewCache(time.Second)
	if c == nil {
		t.Fatalf("expected non-nil cache")
	}
	if c.Entries == nil {
		t.Fatalf("expected Entries map to be initialized")
	}
	if len(c.Entries) != 0 {
		t.Fatalf("expected empty cache, got %d entries", len(c.Entries))
	}

	// Smoke test: add/get works on a fresh cache.
	key := "cache:test:k3"
	val := []byte("value")
	c.Add(key, val)
	got, ok := c.Get(key)
	if !ok || string(got) != string(val) {
		t.Fatalf("add/get failed on fresh cache; ok=%v got=%q want=%q", ok, string(got), string(val))
	}
}

// Ensure two caches from NewCache are independent.
func TestNewCacheIndependentInstances(t *testing.T) {
	c1 := NewCache(time.Second)
	c2 := NewCache(time.Second)

	c1.Add("cache:test:unique", []byte("1"))

	if _, ok := c2.Get("cache:test:unique"); ok {
		t.Fatalf("expected caches to be independent; key from c1 found in c2")
	}
}

// entriesSetForTest is a small helper to set an entry with a specific timestamp.
func (c *Cache) entriesSetForTest(key string, val []byte, createdAt time.Time) {
	c.Entries[key] = cacheEntry{
		createdAt: createdAt,
		val:       val,
	}
}
