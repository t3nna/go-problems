# Observer

Implement a concurrency-safe `EventBus`. Observers subscribe to one topic and receive matching events through `OnEvent`. `Subscribe` returns an unsubscribe function that may safely be called more than once.

## Public API and example

Use `NewEventBus`, `Subscribe(topic, observer)`, and `Publish(Event)`. An observer subscribed to `orders` receives `Event{Topic: "orders"}` but not an event for `payments`.

Publish notifications in subscription order. Concurrent calls to `Subscribe`, unsubscribe, and `Publish` must be safe. Do not hold the bus lock while calling observer code: an observer can be slow or can call back into the bus.

## Acceptance criteria

- Events are delivered only to matching topics.
- Subscription order is preserved for one publish.
- Unsubscribe is idempotent.
- The implementation passes the race detector under concurrent use.

## Tags

`Design Patterns` `Observer` `Concurrency`
