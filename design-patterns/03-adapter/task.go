package main

import "errors"

type LegacyProcessor interface {
	Process(amountDollars float64) (string, error)
}
type PaymentGateway interface {
	Charge(cents int) (string, error)
}

func NewPaymentAdapter(LegacyProcessor) PaymentGateway { return todoGateway{} }

type todoGateway struct{}

func (todoGateway) Charge(int) (string, error) {
	return "", errors.New("TODO: adapt payment processor")
}
