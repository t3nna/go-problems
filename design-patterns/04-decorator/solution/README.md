# Decorator: Solution

## Intent

A Decorator adds behavior while preserving the wrapped object's interface. Both wrappers still satisfy `Sender`, so callers can compose them in any order.

## Go mapping

`SenderFunc` turns functions into interface implementations. `WithRetry` and `WithLogging` each return a new `SenderFunc` that closes over the next sender. This is the same shape as Go HTTP middleware and `http.RoundTripper` wrappers.

## Why use it

Use decorators for cross-cutting behavior such as retries, metrics, tracing, rate limiting, logging, or authorization. The original sender remains focused on sending.

## Trade-offs

Wrapper order matters: logging outside retry logs once, while logging inside retry logs each attempt. Keep retry policies explicit; production retry logic usually also needs backoff and error classification. The solution checks a canceled context before the first call to avoid unnecessary work.
