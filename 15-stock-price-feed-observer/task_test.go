package main

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

type priceUpdate struct {
	symbol string
	price  float64
}

type spyObserver struct {
	mu      sync.Mutex
	updates []priceUpdate
}

func (s *spyObserver) Update(symbol string, price float64) {
	s.mu.Lock()
	s.updates = append(s.updates, priceUpdate{
		symbol: symbol,
		price:  price,
	})
	s.mu.Unlock()
}

func (s *spyObserver) Updates() []priceUpdate {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]priceUpdate, len(s.updates))
	copy(out, s.updates)
	return out
}

type countingObserver struct {
	count atomic.Int64
}

func (c *countingObserver) Update(_ string, _ float64) {
	c.count.Add(1)
}

func TestPriceFeedPublishNotifiesSubscribedObserver(t *testing.T) {
	feed := NewPriceFeed()
	observer := &spyObserver{}

	feed.Subscribe("AAPL", observer)
	feed.Publish("AAPL", 182.3)

	want := []priceUpdate{{symbol: "AAPL", price: 182.3}}
	if got := observer.Updates(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected updates %v, got %v", want, got)
	}
}

func TestPriceFeedPublishNotifiesOnlyMatchingSymbol(t *testing.T) {
	feed := NewPriceFeed()
	aaplObserver := &spyObserver{}
	tslaObserver := &spyObserver{}

	feed.Subscribe("AAPL", aaplObserver)
	feed.Subscribe("TSLA", tslaObserver)

	feed.Publish("AAPL", 200)

	if got := len(aaplObserver.Updates()); got != 1 {
		t.Errorf("Expected AAPL observer to receive 1 update, got %d", got)
	}
	if got := len(tslaObserver.Updates()); got != 0 {
		t.Errorf("Expected TSLA observer to receive 0 updates, got %d", got)
	}
}

func TestPriceFeedUnsubscribeStopsNotifications(t *testing.T) {
	feed := NewPriceFeed()
	observer := &spyObserver{}

	unsubscribe := feed.Subscribe("AAPL", observer)

	feed.Publish("AAPL", 100)
	unsubscribe()
	unsubscribe()
	feed.Publish("AAPL", 101)

	want := []priceUpdate{{symbol: "AAPL", price: 100}}
	if got := observer.Updates(); !reflect.DeepEqual(got, want) {
		t.Errorf("Expected updates %v, got %v", want, got)
	}
}

func TestPriceFeedPublishPreservesSubscriptionOrder(t *testing.T) {
	feed := NewPriceFeed()

	var (
		mu    sync.Mutex
		order []int
	)

	newObserver := func(id int) Observer {
		return observerFunc(func(_ string, _ float64) {
			mu.Lock()
			order = append(order, id)
			mu.Unlock()
		})
	}

	feed.Subscribe("AAPL", newObserver(1))
	feed.Subscribe("AAPL", newObserver(2))
	feed.Subscribe("AAPL", newObserver(3))

	feed.Publish("AAPL", 300)

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(order, want) {
		t.Errorf("Expected notify order %v, got %v", want, order)
	}
}

func TestPriceFeedConcurrentSubscribeAndPublish(t *testing.T) {
	feed := NewPriceFeed()
	observer := &countingObserver{}

	const workers = 50
	const publishes = 20

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsubscribe := feed.Subscribe("AAPL", observer)
			for i := range publishes {
				feed.Publish("AAPL", float64(i))
			}
			unsubscribe()
		}()
	}

	wg.Wait()

	if got := observer.count.Load(); got == 0 {
		t.Fatal("Expected observer to receive updates during concurrent access")
	}
}

type observerFunc func(symbol string, price float64)

func (f observerFunc) Update(symbol string, price float64) {
	f(symbol, price)
}
