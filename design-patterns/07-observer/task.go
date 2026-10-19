package main

type Event struct {
	Topic   string
	Payload any
}
type Observer interface{ OnEvent(Event) }
type EventBus struct{}

func NewEventBus() *EventBus                          { return &EventBus{} }
func (b *EventBus) Subscribe(string, Observer) func() { return func() {} }
func (b *EventBus) Publish(Event)                     {}
