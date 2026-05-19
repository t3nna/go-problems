# Stock Price Feed Observer
Implement a stock price feed using the Observer pattern.

Requirements:

- `NewPriceFeed()` creates an empty feed.
- `Subscribe(symbol string, observer Observer) func()` registers an observer for a stock symbol and returns an unsubscribe function.
- `Publish(symbol string, price float64)` notifies all observers subscribed to that symbol.
- Notifications must be delivered in subscription order.
- Calling unsubscribe multiple times must be safe.
- The feed must be safe for concurrent `Subscribe` and `Publish` calls.

## Tags
`Design Patterns`
