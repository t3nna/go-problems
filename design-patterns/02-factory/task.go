package main

import "errors"

type Channel string

const (
	EmailChannel Channel = "email"
	SMSChannel   Channel = "sms"
)

type Notification struct{ Recipient, Body string }
type Notifier interface {
	Send(Notification) (string, error)
}

func NewNotifier(Channel) (Notifier, error) { return nil, errors.New("TODO: choose notifier") }
