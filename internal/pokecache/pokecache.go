package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries map[string]cacheEntry
	Mutex   sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		Entries: make(map[string]cacheEntry),
	}

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.reapLoop(interval)
		}
	}()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	if _, exist := c.Entries[key]; !exist {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()

		entry := cacheEntry{
			createdAt: time.Now(),
			val:       val,
		}
		c.Entries[key] = entry
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	entry, exist := c.Entries[key]
	if !exist {
		return []byte{}, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for key, value := range c.Entries {
		if time.Since(value.createdAt) >= interval {
			delete(c.Entries, key)
		}
	}
}
