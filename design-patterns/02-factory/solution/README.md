# Factory: Solution

## Intent

A Factory centralizes the choice of concrete implementation while exposing a stable abstraction. Here, callers ask for a `Notifier` and never instantiate an email or SMS implementation directly.

## Go mapping

`NewNotifier` switches on the small `Channel` value and returns an interface. The concrete notifier types are unexported, which keeps their representation private and lets the factory change it later.

## Why use it

Use a factory when configuration, an environment, or a runtime value determines which implementation should serve an interface. It keeps that decision out of business logic and makes each implementation independently testable.

## Trade-offs

For two stable variants a switch is clear. Avoid creating a factory merely to hide a single concrete type; a constructor is simpler. If variants become plugin-like, a registration-based factory may be a better next step.
