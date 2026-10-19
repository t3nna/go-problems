# Builder with Functional Options

Build a validated server configuration without a long, fragile constructor argument list.

Implement `NewServerConfig(opts ...Option)`. It must start with these defaults: host `localhost`, port `8080`, timeout `30s`, TLS disabled, and 100 maximum connections. Apply options in order and return an error for an empty host, a port outside `1..65535`, a non-positive timeout, or non-positive connection limit.

The public API is `ServerConfig`, `Option`, `NewServerConfig`, and the five `With...` helpers in `task.go`.

For example, `NewServerConfig(WithPort(8443), WithTLS(true))` keeps the default host and timeout while changing only those two settings.

## Acceptance criteria

- Defaults are useful when no options are passed.
- Every option overrides only its own field.
- Invalid input returns an error rather than creating a partial configuration.

## Tags

`Design Patterns` `Builder` `Functional Options`
