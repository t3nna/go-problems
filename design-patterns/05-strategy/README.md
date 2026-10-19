# Strategy

Shipping rules change independently of checkout. Implement `NewCheckout(strategy ShippingStrategy)` and `(*Checkout).Total(order Order)` so checkout delegates the shipping calculation to the selected strategy.

`FlatRate(cents)` always charges that amount. `WeightBased(baseCents, perStartedKilogramCents)` charges the base plus the per-kilogram amount for every started kilogram; 1001 grams therefore counts as two kilograms. `Total` adds shipping to a non-negative subtotal.

## Public API and example

Keep the supplied `Order`, `ShippingStrategy`, `Checkout`, `NewCheckout`, `FlatRate`, and `WeightBased` API. For example, a 1,500-cent order using `FlatRate(500)` totals 2,000 cents.

## Acceptance criteria

- Checkout works with either strategy without conditional logic in `Total`.
- A nil strategy and negative subtotal return errors.
- Strategy configuration and invalid order values are validated.

## Tags

`Design Patterns` `Strategy`
