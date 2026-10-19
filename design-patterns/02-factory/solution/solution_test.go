package main

import "testing"

func TestNewNotifierCreatesEmailNotifier(t *testing.T) {
	n, err := NewNotifier(EmailChannel)
	if err != nil {
		t.Fatal(err)
	}
	got, err := n.Send(Notification{"ada@example.com", "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "email:ada@example.com:hello" {
		t.Fatalf("got %q", got)
	}
}
func TestNewNotifierCreatesSMSNotifier(t *testing.T) {
	n, err := NewNotifier(SMSChannel)
	if err != nil {
		t.Fatal(err)
	}
	got, err := n.Send(Notification{"+48123456789", "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "sms:+48123456789:hello" {
		t.Fatalf("got %q", got)
	}
}
func TestNewNotifierRejectsUnknownChannel(t *testing.T) {
	if _, err := NewNotifier("push"); err == nil {
		t.Fatal("expected error")
	}
}
func TestNotifierRejectsIncompleteNotification(t *testing.T) {
	n, _ := NewNotifier(EmailChannel)
	if _, err := n.Send(Notification{}); err == nil {
		t.Fatal("expected error")
	}
}
