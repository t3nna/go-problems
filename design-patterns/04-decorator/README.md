# Decorator

Add retry and logging behavior around a `Sender` without changing its interface or modifying its concrete implementation.

Implement `WithRetry(next Sender, attempts int)` and `WithLogging(next Sender, logf func(string))`. `attempts` is the total number of sends, including the first. Retrying must stop on success, return the final send error, and avoid calling the wrapped sender when the context is already canceled. Logging must emit `sending to=<recipient>` and then either `sent to=<recipient>` or `failed to=<recipient>`.

## Public API and example

Both functions return a `Sender`, so they can be composed: `WithLogging(WithRetry(sender, 3), logf)`. Do not change the `Sender` interface or `Message` type.

## Acceptance criteria

- Retry calls the wrapped sender at most `attempts` times.
- Logging preserves the wrapped sender's result.
- A canceled context stops retrying before the first wrapped call.

## Tags

`Design Patterns` `Decorator` `Composition`
