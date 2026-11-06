package main

import (
	"sync"
	"time"
)

type TtlCache struct {
	data map[string]*entry
	heap minHeap
	mu   sync.Mutex
	stop chan struct{}
	done chan struct{}
}

func NewTtlCache() *TtlCache {
	c := &TtlCache{
		data: make(map[string]*entry),
	}

	go c.clear()
	return c

}

func (c *TtlCache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	oldEntry, ok := c.data[key]
	if ok && oldEntry.index >= 0 {
		// remove because it can have other value + ttl which changes postion in heap
		c.heap.remove(oldEntry.index)
	}

	e := &entry{
		key:   key,
		value: value,
		index: -1, // by default not expiring
	}

	if ttl > 0 {
		e.expiresAt = time.Now().Add(ttl)
		c.heap.insert(e)
	}

	c.data[key] = e
}

func (c *TtlCache) Get(key string) (string, bool) {
	e, ok := c.data[key]
	if !ok || (!e.expiresAt.IsZero() && !time.Now().Before(e.expiresAt)) {
		return "", false
	}
	return e.value, true
}

func (c *TtlCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.data[key]
	if !ok {
		return
	}
	c.heap.remove(e.index)
	delete(c.data, key)

}

func (c *TtlCache) Stop() {
	close(c.stop)
	<-c.done
}

func (c *TtlCache) clear() {
	defer close(c.done)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stop:
			return
		}
	}

}

func (c *TtlCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for c.heap.length() < 0 {
		e := c.heap.peek()

		if time.Now().Before(e.expiresAt) {
			return
		}
		c.Delete(e.key)
	}
}
