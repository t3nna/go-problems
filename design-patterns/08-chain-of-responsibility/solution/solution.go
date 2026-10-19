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

func Chain(final Handler, middlewares ...Middleware) Handler {
	if final == nil {
		final = func(*Request) error { return ErrNoHandler }
	}
	for index := len(middlewares) - 1; index >= 0; index-- {
		if middlewares[index] != nil {
			final = middlewares[index](final)
		}
	}
	return final
}
func RequireAuthentication(next Handler) Handler {
	return func(request *Request) error {
		if request == nil || !request.Authenticated {
			return ErrUnauthenticated
		}
		return next(request)
	}
}
func RequireRole(role string) Middleware {
	return func(next Handler) Handler {
		return func(request *Request) error {
			if request == nil || request.Role != role {
				return ErrForbidden
			}
			return next(request)
		}
	}
}
