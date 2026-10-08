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
	ent := c.cache[address]

	if ent == nil {
		t := &task{
			isReady: make(chan struct{}),
		}
		c.cache[address] = t
		c.mu.Unlock()

		v, err := c.client.Get(address)

		t.body = v
		t.err = err

		close(t.isReady)

		return v, err
	} else {

		<-ent.isReady
		c.mu.Unlock()
		return ent.body, ent.err
	}

}
