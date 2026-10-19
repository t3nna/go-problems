# Strategy: Solution

## Intent

The Strategy pattern lets one part of an algorithm vary independently. Here, checkout always calculates `subtotal + shipping`, while the shipping rule is selected by the caller.

## Go mapping

`ShippingStrategy` is a one-method interface. `FlatRate` and `FreeShipping` are ordinary small values that implement it. `Checkout` stores the interface and calls `ShippingCost`; it does not check which concrete strategy it received.

```go
checkout := NewCheckout(FreeShipping{})
total := checkout.Total(1500) // 1500
```

The custom-strategy test is the important part: a new type can implement `ShippingCost` and work with checkout without changing checkout code.

## When to use it

Use Strategy when a policy varies: discounts, tax rules, sorting, authorization, retry decisions, or serialization choices. The caller chooses the policy, while the context owns the workflow that uses it.

## Trade-offs

This deliberately small example assumes valid subtotals and a non-nil strategy so the pattern stays visible. In production, validate those values at the API boundary. Do not introduce a strategy interface for a policy that will never vary; a direct function is clearer in that case.
