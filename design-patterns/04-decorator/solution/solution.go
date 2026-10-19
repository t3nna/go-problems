package main

import (
	"context"
	"fmt"
)

type Message struct{ To, Body string }
type Sender interface {
	Send(context.Context, Message) error
}
type SenderFunc func(context.Context, Message) error

func (f SenderFunc) Send(ctx context.Context, message Message) error { return f(ctx, message) }

func WithRetry(next Sender, attempts int) Sender {
	return SenderFunc(func(ctx context.Context, message Message) error {
		if next == nil {
			return fmt.Errorf("sender is required")
		}
		if attempts <= 0 {
			return fmt.Errorf("attempts must be positive")
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		var err error
		for range attempts {
			if err = next.Send(ctx, message); err == nil {
				return nil
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}
		return err
	})
}
func WithLogging(next Sender, logf func(string)) Sender {
	return SenderFunc(func(ctx context.Context, message Message) error {
		if next == nil {
			return fmt.Errorf("sender is required")
		}
		if logf != nil {
			logf("sending to=" + message.To)
		}
		err := next.Send(ctx, message)
		if logf != nil {
			if err == nil {
				logf("sent to=" + message.To)
			} else {
				logf("failed to=" + message.To)
			}
		}
		return err
	})
}
