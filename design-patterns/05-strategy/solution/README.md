# Strategy: Solution

## Intent

The Strategy pattern defines a family of algorithms behind one interface and lets the caller select one. Checkout owns the stable workflow; a `ShippingStrategy` owns the variable shipping rule.

## Go mapping

The solution uses a small interface and the `shippingStrategyFunc` adapter, allowing constructor functions to return strategies without exported concrete types. `Checkout.Total` has no knowledge of flat or weight-based pricing.

## Why use it

Use Strategy when one policy changes independently from the code that invokes it: pricing, authorization, serialization, sorting, retry rules, and feature rollout logic are common examples.

## Trade-offs

Strategies add an abstraction and validation boundary. They are worthwhile when variations are real or expected, not simply to avoid one small conditional. The explicit constructor functions keep the public API easy to read at the call site.
