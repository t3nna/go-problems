# Strategy

Checkout always adds shipping to the order subtotal, but the shipping rule can change. Implement checkout so it depends on a `ShippingStrategy`, not on a specific shipping rule.

This exercise intentionally uses two very simple strategies:

- `FlatRate{Cents: 500}` always charges 500 cents.
- `FreeShipping{}` always charges 0 cents.

## Public API and example

Keep `ShippingStrategy`, `FlatRate`, `FreeShipping`, `Checkout`, and `NewCheckout` as defined in `task.go`.

```go
checkout := NewCheckout(FlatRate{Cents: 500})
total := checkout.Total(1500) // 2000
```

The checkout code must work unchanged with `FreeShipping{}` or any other value that implements `ShippingStrategy`.

## Acceptance criteria

- `Total` adds the selected strategy's shipping cost to the subtotal.
- `FlatRate` returns its configured number of cents.
- `FreeShipping` returns zero.
- A caller can add a new strategy without editing `Checkout`.

## Tags

`Design Patterns` `Strategy`
