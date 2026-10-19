# Chain of Responsibility

Compose request-processing middleware. Implement `Chain(final Handler, middlewares ...Middleware)` so middleware run in the order declared, wrapping the final handler. A middleware may return an error without calling `next`, which stops the chain.

Implement `RequireAuthentication` and `RequireRole`. Authentication must fail with `ErrUnauthenticated`; role mismatches fail with `ErrForbidden`. A nil final handler must return `ErrNoHandler`.

## Public API and example

Keep `Handler`, `Middleware`, `Chain`, `RequireAuthentication`, and `RequireRole`. For example, `Chain(final, RequireAuthentication, RequireRole("admin"))` must reject an unauthenticated request before the role check or final handler runs.

## Acceptance criteria

- Middleware run in the same order passed to `Chain`.
- Returning an error prevents subsequent handlers from running.
- Nil final handlers return `ErrNoHandler`.

## Tags

`Design Patterns` `Chain of Responsibility` `Middleware`
