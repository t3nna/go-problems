package main

import (
	"sync"
	"time"
)

type entry struct {
	key       string
	value     string
	expiresAt time.Time // Zero means no expiration.
	index     int       // Index in the heap, or -1 if not in the heap.
}

// minHeap stores entries ordered by their expiration time.
type minHeap struct {
	data []*entry
}

func (h *minHeap) length() int {
	return len(h.data)
}

func (h *minHeap) peek() *entry {
	if h.length() == 0 {
		return nil
	}
	return h.data[0]
}

func (h *minHeap) insert(e *entry) {
	e.index = h.length()
	h.data = append(h.data, e)
	h.heapifyUp(e.index)
}

// remove deletes and returns the entry at index.
func (h *minHeap) remove(index int) *entry {
	last := h.length() - 1
	removed := h.data[index]

	h.swap(index, last)
	h.data[last] = nil
	h.data = h.data[:last]
	removed.index = -1

	if index < h.length() {
		parent := (index - 1) / 2
		if index > 0 && h.less(index, parent) {
			h.heapifyUp(index)
		} else {
			h.heapifyDown(index)
		}
	}

	return removed
}

func (h *minHeap) deleteMin() *entry {
	if h.length() == 0 {
		return nil
	}
	return h.remove(0)
}

func (h *minHeap) heapifyUp(index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if !h.less(index, parent) {
			return
		}

		h.swap(index, parent)
		index = parent
	}
}

func (h *minHeap) heapifyDown(index int) {
	for {
		left := 2*index + 1
		right := 2*index + 2
		smallest := index

		if left < h.length() && h.less(left, smallest) {
			smallest = left
		}
		if right < h.length() && h.less(right, smallest) {
			smallest = right
		}
		if smallest == index {
			return
		}

		h.swap(index, smallest)
		index = smallest
	}
}

func (h *minHeap) less(i, j int) bool {
	return h.data[i].expiresAt.Before(h.data[j].expiresAt)
}

func (h *minHeap) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.data[i].index = i
	h.data[j].index = j
}

type TtlCache struct {
	mu      sync.RWMutex
	entries map[string]*entry
	expires minHeap

	stopOnce sync.Once
	stop     chan struct{}
	done     chan struct{}
}

func NewTtlCache() *TtlCache {
	c := &TtlCache{
		entries: make(map[string]*entry),
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
	go c.clear()
	return c
}

// Set stores a value. A TTL <= 0 means the value does not expire.
func (c *TtlCache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if old, ok := c.entries[key]; ok && old.index >= 0 {
		c.expires.remove(old.index)
	}

	e := &entry{
		key:   key,
		value: value,
		index: -1,
	}

	if ttl > 0 {
		e.expiresAt = time.Now().Add(ttl)
		c.expires.insert(e)
	}

	c.entries[key] = e
}

func (c *TtlCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	e, ok := c.entries[key]
	if !ok || (!e.expiresAt.IsZero() && !time.Now().Before(e.expiresAt)) {
		return "", false
	}
	return e.value, true
}

func (c *TtlCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.entries[key]
	if !ok {
		return
	}

	if e.index >= 0 {
		c.expires.remove(e.index)
	}
	delete(c.entries, key)
}

// Stop stops background cleanup and waits for it to finish.
// It is safe to call more than once.
func (c *TtlCache) Stop() {
	c.stopOnce.Do(func() { close(c.stop) })
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

	now := time.Now()
	for c.expires.length() > 0 {
		e := c.expires.peek()
		if now.Before(e.expiresAt) {
			break
		}

		c.expires.deleteMin()
		delete(c.entries, e.key)
	}
}
