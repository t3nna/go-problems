package main

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

type eventSpy struct {
	mu     sync.Mutex
	events []Event
}

func (s *eventSpy) OnEvent(event Event) {
	s.mu.Lock()
	s.events = append(s.events, event)
	s.mu.Unlock()
}
func (s *eventSpy) Events() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Event(nil), s.events...)
}

type countingObserver struct{ count atomic.Int64 }

func (o *countingObserver) OnEvent(Event) { o.count.Add(1) }

func TestEventBusNotifiesOnlyTopicSubscribers(t *testing.T) {
	bus := NewEventBus()
	orders, payments := &eventSpy{}, &eventSpy{}
	bus.Subscribe("orders", orders)
	bus.Subscribe("payments", payments)
	event := Event{Topic: "orders", Payload: 42}
	bus.Publish(event)
	if !reflect.DeepEqual(orders.Events(), []Event{event}) {
		t.Fatalf("orders = %#v", orders.Events())
	}
	if len(payments.Events()) != 0 {
		t.Fatalf("payments = %#v", payments.Events())
	}
}
func TestEventBusPreservesSubscriptionOrder(t *testing.T) {
	bus := NewEventBus()
	var got []int
	var mu sync.Mutex
	observer := func(id int) Observer {
		return observerFunc(func(Event) { mu.Lock(); got = append(got, id); mu.Unlock() })
	}
	bus.Subscribe("topic", observer(1))
	bus.Subscribe("topic", observer(2))
	bus.Subscribe("topic", observer(3))
	bus.Publish(Event{Topic: "topic"})
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("order = %v", got)
	}
}
func TestEventBusUnsubscribeIsIdempotent(t *testing.T) {
	bus := NewEventBus()
	observer := &eventSpy{}
	unsubscribe := bus.Subscribe("topic", observer)
	bus.Publish(Event{Topic: "topic"})
	unsubscribe()
	unsubscribe()
	bus.Publish(Event{Topic: "topic"})
	if got := len(observer.Events()); got != 1 {
		t.Fatalf("events = %d", got)
	}
}
func TestEventBusSupportsConcurrentPublishAndSubscribe(t *testing.T) {
	bus := NewEventBus()
	observer := &countingObserver{}
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsubscribe := bus.Subscribe("topic", observer)
			for range 10 {
				bus.Publish(Event{Topic: "topic"})
			}
			unsubscribe()
		}()
	}
	wg.Wait()
	if observer.count.Load() == 0 {
		t.Fatal("expected observed events")
	}
}

type observerFunc func(Event)

func (f observerFunc) OnEvent(event Event) { f(event) }
