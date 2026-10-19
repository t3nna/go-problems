package main

import (
	"context"
	"errors"
	"testing"
)

func TestWithRetryRetriesUntilSuccess(t *testing.T) {
	calls := 0
	sender := SenderFunc(func(context.Context, Message) error {
		calls++
		if calls < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if err := WithRetry(sender, 3).Send(context.Background(), Message{}); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d", calls)
	}
}
func TestWithRetryReturnsLastError(t *testing.T) {
	want := errors.New("failed")
	calls := 0
	sender := SenderFunc(func(context.Context, Message) error { calls++; return want })
	if err := WithRetry(sender, 2).Send(context.Background(), Message{}); !errors.Is(err, want) || calls != 2 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}
func TestWithLoggingPreservesSenderAndReportsOutcome(t *testing.T) {
	var logs []string
	sender := SenderFunc(func(context.Context, Message) error { return nil })
	if err := WithLogging(sender, func(line string) { logs = append(logs, line) }).Send(context.Background(), Message{To: "ada"}); err != nil {
		t.Fatal(err)
	}
	if len(logs) != 2 || logs[0] != "sending to=ada" || logs[1] != "sent to=ada" {
		t.Fatalf("logs = %v", logs)
	}
}
func TestWithRetryHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	sender := SenderFunc(func(context.Context, Message) error { calls++; return nil })
	if err := WithRetry(sender, 2).Send(ctx, Message{}); !errors.Is(err, context.Canceled) || calls != 0 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}
