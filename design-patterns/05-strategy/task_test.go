package main

import "testing"

func TestCheckoutUsesFlatRateStrategy(t *testing.T) {
	checkout := NewCheckout(FlatRate{Cents: 500})

	if got := checkout.Total(1500); got != 2000 {
		t.Fatalf("Total() = %d, want 2000", got)
	}
}

func TestCheckoutUsesFreeShippingStrategy(t *testing.T) {
	checkout := NewCheckout(FreeShipping{})

	if got := checkout.Total(1500); got != 1500 {
		t.Fatalf("Total() = %d, want 1500", got)
	}
}

func TestCheckoutUsesCustomStrategy(t *testing.T) {
	checkout := NewCheckout(shippingCost(250))

	if got := checkout.Total(1500); got != 1750 {
		t.Fatalf("Total() = %d, want 1750", got)
	}
}

type shippingCost int

func (cost shippingCost) ShippingCost() int { return int(cost) }
