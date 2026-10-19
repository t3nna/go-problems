package main

import "testing"

func TestNewNotifierCreatesEmailNotifier(t *testing.T) {
	notifier, err := NewNotifier(EmailChannel)
	if err != nil {
		t.Fatal(err)
	}
	got, err := notifier.Send(Notification{Recipient: "ada@example.com", Body: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "email:ada@example.com:hello" {
		t.Fatalf("Send() = %q", got)
	}
}
func TestNewNotifierCreatesSMSNotifier(t *testing.T) {
	notifier, err := NewNotifier(SMSChannel)
	if err != nil {
		t.Fatal(err)
	}
	got, err := notifier.Send(Notification{Recipient: "+48123456789", Body: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "sms:+48123456789:hello" {
		t.Fatalf("Send() = %q", got)
	}
}
func TestNewNotifierRejectsUnknownChannel(t *testing.T) {
	if _, err := NewNotifier("push"); err == nil {
		t.Fatal("expected unsupported channel error")
	}
}
func TestNotifierRejectsIncompleteNotification(t *testing.T) {
	notifier, _ := NewNotifier(EmailChannel)
	if _, err := notifier.Send(Notification{}); err == nil {
		t.Fatal("expected validation error")
	}
}
