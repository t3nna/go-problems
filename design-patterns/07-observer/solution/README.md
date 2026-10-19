# Observer: Solution

## Intent

Observer creates a one-to-many dependency: when a subject changes, interested observers receive a notification. The generic event bus decouples publishers from subscribers and filters by topic.

## Go mapping

`Observer` is a tiny interface, while `EventBus` owns subscriptions. The bus copies matching observers under a read lock and releases the lock before invoking callbacks. `sync.Once` makes the returned unsubscribe function idempotent.

## Why use it

Use Observer for in-process events, cache invalidation, UI updates, monitoring hooks, and domain notifications. Publishers only know an event topic, not concrete listeners.

## Trade-offs

Observer code can hide control flow and make callback latency surprising. Copying before notification avoids lock reentrancy and long lock holds, but an unsubscribe concurrent with publishing may still receive the already-snapshotted event. For cross-process delivery, use a durable message broker rather than this in-memory pattern.
