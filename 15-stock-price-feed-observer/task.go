package main

import "sync"

type Observer interface {
	Update(symbol string, price float64)
}

type subscriber struct {
	id       int
	observer Observer
}

type PriceFeed struct {
	mu          sync.RWMutex
	nextID      int
	subscribers map[string][]subscriber
}

func NewPriceFeed() *PriceFeed {
	return &PriceFeed{
		subscribers: make(map[string][]subscriber),
	}
}

func (pf *PriceFeed) Subscribe(symbol string, observer Observer) func() {
	if observer == nil {
		return func() {}
	}

	pf.mu.Lock()
	id := pf.nextID
	pf.nextID++
	pf.subscribers[symbol] = append(pf.subscribers[symbol], subscriber{
		id:       id,
		observer: observer,
	})
	pf.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			pf.remove(symbol, id)
		})
	}
}

func (pf *PriceFeed) Publish(symbol string, price float64) {
	pf.mu.RLock()
	list := pf.subscribers[symbol]
	observers := make([]Observer, 0, len(list))
	for _, sub := range list {
		observers = append(observers, sub.observer)
	}
	pf.mu.RUnlock()

	for _, observer := range observers {
		observer.Update(symbol, price)
	}
}

func (pf *PriceFeed) remove(symbol string, id int) {
	pf.mu.Lock()
	defer pf.mu.Unlock()

	list := pf.subscribers[symbol]
	for i := range list {
		if list[i].id != id {
			continue
		}

		pf.subscribers[symbol] = append(list[:i], list[i+1:]...)
		if len(pf.subscribers[symbol]) == 0 {
			delete(pf.subscribers, symbol)
		}
		return
	}
}
