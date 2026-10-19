# Adapter

The application charges integer cents, but a legacy payment library accepts floating-point dollars. Implement `NewPaymentAdapter(legacy LegacyProcessor) PaymentGateway` so that callers only use `Charge(cents int)`.

Convert cents to dollars before delegating. Reject non-positive charges and a nil legacy processor, and preserve errors returned by the legacy processor.

## Public API

Implement `NewPaymentAdapter(LegacyProcessor) PaymentGateway`. Calling `Charge(1234)` must invoke `LegacyProcessor.Process(12.34)` and return its receipt or error.

## Acceptance criteria

- Charging 1234 cents delegates exactly 12.34 dollars.
- The adapter keeps the legacy dependency out of the gateway-facing API.
- Invalid input and legacy failures are returned as errors.

## Tags

`Design Patterns` `Adapter`
