package main

import "testing"

func TestCheckoutUsesFlatRateStrategy(t *testing.T) {
	checkout, err := NewCheckout(FlatRate(500))
	if err != nil {
		t.Fatal(err)
	}
	got, err := checkout.Total(Order{SubtotalCents: 1500, WeightGrams: 1})
	if err != nil || got != 2000 {
		t.Fatalf("total=%d err=%v", got, err)
	}
}
func TestCheckoutUsesWeightBasedStrategy(t *testing.T) {
	checkout, err := NewCheckout(WeightBased(200, 75))
	if err != nil {
		t.Fatal(err)
	}
	got, err := checkout.Total(Order{SubtotalCents: 1000, WeightGrams: 1001})
	if err != nil || got != 1350 {
		t.Fatalf("total=%d err=%v", got, err)
	}
}
func TestCheckoutRejectsInvalidInputs(t *testing.T) {
	if _, err := NewCheckout(nil); err == nil {
		t.Fatal("expected strategy error")
	}
	checkout, _ := NewCheckout(FlatRate(100))
	if _, err := checkout.Total(Order{SubtotalCents: -1}); err == nil {
		t.Fatal("expected order error")
	}
}
