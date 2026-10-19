package main

import "errors"

var (
	ErrNoHandler       = errors.New("no handler")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")
)

type Request struct {
	User, Role    string
	Authenticated bool
}
type Handler func(*Request) error
type Middleware func(Handler) Handler

func Chain(Handler, ...Middleware) Handler { return func(*Request) error { return ErrNoHandler } }
func RequireAuthentication(Handler) Handler {
	return func(*Request) error { return ErrUnauthenticated }
}
func RequireRole(string) Middleware {
	return func(Handler) Handler { return func(*Request) error { return ErrForbidden } }
}
