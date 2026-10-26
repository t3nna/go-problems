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
	return &Cache{
		cache:  make(map[string]*task),
		client: client,
	}
}

// Cache Client.Get result
func (c *Cache) Get(address string) (string, error) {
	c.mu.Lock()
	t := c.cache[address]

	if t == nil {
		t = &task{isReady: make(chan struct{})}
		c.cache[address] = t

		res, err := c.client.Get(address)

		t.body = res
		t.err = err
		c.mu.Unlock()

		close(t.isReady)

		return res, err
	}
	c.mu.Unlock()

	<-t.isReady

	return t.body, t.err
}
