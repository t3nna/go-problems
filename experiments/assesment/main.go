package main

import "time"

type Item struct {
	data      string
	liveUntil time.Time
	ttl       int
}

type Cache struct {
	c map[string]*Item
}

func (c *Cache) Get(id string) {

}
func (c *Cache) Set(id, value string, ttl int) {
	c.c[id] = &Item{
		data:      value,
		liveUntil: time.Now() + (time.Second * ttl),
		ttl:       ttl,
	}

}

func main() {
}
