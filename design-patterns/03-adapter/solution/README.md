# Adapter: Solution

## Intent

An Adapter translates one interface into another expected by the client. It lets new code use the more suitable `PaymentGateway` contract while preserving a legacy dependency.

## Go mapping

`paymentAdapter` stores the `LegacyProcessor` and implements `PaymentGateway`. Its `Charge` method performs unit conversion and delegates the work. The concrete adapter is private; callers obtain it through `NewPaymentAdapter`.

## Why use it

Use an adapter at a system boundary: a third-party SDK, a legacy package, or a transport layer with an inconvenient API. It localizes conversion, validation, and error behavior.

## Trade-offs

The adapter must make conversion rules explicit. Floating-point currency is normally unsuitable for internal domain logic, which is why the public API retains integer cents. Avoid an adapter when you control both interfaces and can simplify them directly.
