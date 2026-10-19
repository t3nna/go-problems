package main

import (
	"errors"
	"testing"
)

type processorSpy struct {
	dollars float64
	err     error
}

func (p *processorSpy) Process(amount float64) (string, error) {
	p.dollars = amount
	if p.err != nil {
		return "", p.err
	}
	return "receipt", nil
}

func TestPaymentAdapterConvertsCentsToDollars(t *testing.T) {
	legacy := &processorSpy{}
	receipt, err := NewPaymentAdapter(legacy).Charge(1234)
	if err != nil {
		t.Fatal(err)
	}
	if receipt != "receipt" || legacy.dollars != 12.34 {
		t.Fatalf("receipt=%q dollars=%v", receipt, legacy.dollars)
	}
}
func TestPaymentAdapterPropagatesLegacyError(t *testing.T) {
	want := errors.New("declined")
	_, err := NewPaymentAdapter(&processorSpy{err: want}).Charge(100)
	if !errors.Is(err, want) {
		t.Fatalf("error = %v", err)
	}
}
func TestPaymentAdapterRejectsInvalidCharge(t *testing.T) {
	if _, err := NewPaymentAdapter(&processorSpy{}).Charge(0); err == nil {
		t.Fatal("expected error")
	}
}
func TestPaymentAdapterHandlesNilLegacyProcessor(t *testing.T) {
	if _, err := NewPaymentAdapter(nil).Charge(100); err == nil {
		t.Fatal("expected error")
	}
}
