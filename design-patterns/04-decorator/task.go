package main

import (
	"context"
	"errors"
)

type Message struct{ To, Body string }
type Sender interface {
	Send(context.Context, Message) error
}
type SenderFunc func(context.Context, Message) error

func (f SenderFunc) Send(ctx context.Context, message Message) error { return f(ctx, message) }

func WithRetry(Sender, int) Sender            { return todoSender{} }
func WithLogging(Sender, func(string)) Sender { return todoSender{} }

type todoSender struct{}

func (todoSender) Send(context.Context, Message) error { return errors.New("TODO: decorate sender") }
