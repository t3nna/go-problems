package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestChainRunsMiddlewareInDeclaredOrder(t *testing.T) {
	var calls []string
	middleware := func(name string) Middleware {
		return func(next Handler) Handler {
			return func(request *Request) error {
				calls = append(calls, name+" before")
				err := next(request)
				calls = append(calls, name+" after")
				return err
			}
		}
	}
	handler := Chain(func(*Request) error { calls = append(calls, "final"); return nil }, middleware("first"), middleware("second"))
	if err := handler(&Request{}); err != nil {
		t.Fatal(err)
	}
	want := []string{"first before", "second before", "final", "second after", "first after"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %v", calls)
	}
}
func TestRequireAuthenticationStopsChain(t *testing.T) {
	called := false
	handler := Chain(func(*Request) error { called = true; return nil }, RequireAuthentication)
	if err := handler(&Request{}); !errors.Is(err, ErrUnauthenticated) || called {
		t.Fatalf("err=%v called=%v", err, called)
	}
}
func TestRequireRoleAllowsOnlyMatchingRole(t *testing.T) {
	handler := Chain(func(*Request) error { return nil }, RequireAuthentication, RequireRole("admin"))
	if err := handler(&Request{Authenticated: true, Role: "member"}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("err=%v", err)
	}
	if err := handler(&Request{Authenticated: true, Role: "admin"}); err != nil {
		t.Fatal(err)
	}
}
func TestChainHandlesNilFinalHandler(t *testing.T) {
	if err := Chain(nil)(&Request{}); !errors.Is(err, ErrNoHandler) {
		t.Fatalf("err=%v", err)
	}
}
