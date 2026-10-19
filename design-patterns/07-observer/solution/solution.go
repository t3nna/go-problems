package main

import "sync"

type Event struct {
	Topic   string
	Payload any
}
type Observer interface{ OnEvent(Event) }
type subscriber struct {
	id       int
	observer Observer
}
type EventBus struct {
	mu          sync.RWMutex
	nextID      int
	subscribers map[string][]subscriber
}

func NewEventBus() *EventBus { return &EventBus{subscribers: make(map[string][]subscriber)} }
func (b *EventBus) Subscribe(topic string, observer Observer) func() {
	if observer == nil {
		return func() {}
	}
	b.mu.Lock()
	id := b.nextID
	b.nextID++
	b.subscribers[topic] = append(b.subscribers[topic], subscriber{id: id, observer: observer})
	b.mu.Unlock()
	var once sync.Once
	return func() { once.Do(func() { b.remove(topic, id) }) }
}
func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	subscriptions := b.subscribers[event.Topic]
	observers := make([]Observer, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		observers = append(observers, subscription.observer)
	}
	b.mu.RUnlock()
	for _, observer := range observers {
		observer.OnEvent(event)
	}
}
func (b *EventBus) remove(topic string, id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	subscriptions := b.subscribers[topic]
	for index, subscription := range subscriptions {
		if subscription.id != id {
			continue
		}
		b.subscribers[topic] = append(subscriptions[:index], subscriptions[index+1:]...)
		if len(b.subscribers[topic]) == 0 {
			delete(b.subscribers, topic)
		}
		return
	}
}
