package main

import "fmt"

type LegacyProcessor interface {
	Process(amountDollars float64) (string, error)
}
type PaymentGateway interface {
	Charge(cents int) (string, error)
}

func NewPaymentAdapter(legacy LegacyProcessor) PaymentGateway { return paymentAdapter{legacy: legacy} }

type paymentAdapter struct{ legacy LegacyProcessor }

func (a paymentAdapter) Charge(cents int) (string, error) {
	if a.legacy == nil {
		return "", fmt.Errorf("legacy processor is required")
	}
	if cents <= 0 {
		return "", fmt.Errorf("charge must be positive")
	}
	return a.legacy.Process(float64(cents) / 100)
}
