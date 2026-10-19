# Builder with Functional Options: Solution

## Intent

The Builder pattern separates construction from representation. Go commonly expresses it with functional options: callers start from safe defaults and name only the settings they need.

## Go mapping

`Option` is a function that mutates a private-in-construction `ServerConfig`. `NewServerConfig` owns the initial defaults and applies the options in call order. Each `With...` function validates its input before changing the configuration.

## Why use it

This avoids constructors such as `NewServer(host, port, timeout, tls, maxConnections)`, where same-typed arguments are easy to misorder and optional arguments grow poorly. It is especially useful for exported APIs expected to acquire more optional configuration.

## Trade-offs

Functional options are less appropriate for a tiny value with only one or two mandatory fields. Keep mandatory dependencies as normal parameters, and document defaults because options otherwise hide them.

The solution rejects invalid options before returning a configuration, so callers never receive a partially valid result.
