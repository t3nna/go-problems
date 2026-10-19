# Chain of Responsibility: Solution

## Intent

Chain of Responsibility passes a request through a sequence of handlers. Each handler can process, enrich, reject, or pass the request along without the caller knowing which handler ultimately acts.

## Go mapping

Go middleware naturally models the pattern as `func(Handler) Handler`. `Chain` wraps from right to left so the middleware supplied first receives the request first. Authentication and authorization return early rather than invoking `next` when they reject a request.

## Why use it

Use this pattern for HTTP middleware, validation pipelines, approval flows, logging, metrics, and authorization. Independent responsibilities remain small and composable.

## Trade-offs

Ordering is part of the behavior: role checks must follow authentication, and observability middleware may need to wrap the whole chain. Long chains can make execution paths less obvious, so keep each component focused and document required order.
