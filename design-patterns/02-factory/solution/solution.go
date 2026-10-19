package main

import "fmt"

type Channel string

const (
	EmailChannel Channel = "email"
	SMSChannel   Channel = "sms"
)

type Notification struct{ Recipient, Body string }
type Notifier interface {
	Send(Notification) (string, error)
}

func NewNotifier(channel Channel) (Notifier, error) {
	switch channel {
	case EmailChannel:
		return emailNotifier{}, nil
	case SMSChannel:
		return smsNotifier{}, nil
	default:
		return nil, fmt.Errorf("unsupported channel %q", channel)
	}
}

type emailNotifier struct{}

func (emailNotifier) Send(n Notification) (string, error) {
	if err := validate(n); err != nil {
		return "", err
	}
	return "email:" + n.Recipient + ":" + n.Body, nil
}

type smsNotifier struct{}

func (smsNotifier) Send(n Notification) (string, error) {
	if err := validate(n); err != nil {
		return "", err
	}
	return "sms:" + n.Recipient + ":" + n.Body, nil
}
func validate(n Notification) error {
	if n.Recipient == "" || n.Body == "" {
		return fmt.Errorf("recipient and body are required")
	}
	return nil
}
