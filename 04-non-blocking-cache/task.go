package main

import (
	"sync"
)

type Client interface {
	Get(address string) (string, error)
}

type task struct {
	body    string
	err     error
	isReady chan struct{}
}

type Cache struct {
	client Client
	// You can add new fields if needed
	cache map[string]*task
	mu    sync.Mutex
}

// Don't update signature of NewCache
func NewCache(client Client) *Cache {
	// TODO: Implement
	return &Cache{client: client,
		cache: make(map[string]*task),
	}
}

// Cache Client.Get result
func (c *Cache) Get(address string) (string, error) {
	// TODO: Implement. Right now it doesn't cache

	c.mu.Lock()
	v := c.cache[address]

	if v == nil {
		t := &task{isReady: make(chan struct{})}
		c.cache[address] = t
		c.mu.Unlock()
		v = t
		res, err := c.client.Get(address)
		v.body = res
		v.err = err
		close(v.isReady)
	} else {
		c.mu.Unlock()
		<-v.isReady
	}
	return v.body, v.err
}
