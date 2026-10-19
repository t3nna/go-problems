package main

import "errors"

type Order struct{ SubtotalCents, WeightGrams int }
type ShippingStrategy interface{ ShippingCost(Order) (int, error) }
type Checkout struct{ strategy ShippingStrategy }

func NewCheckout(ShippingStrategy) (*Checkout, error) {
	return &Checkout{}, nil
}
func (c *Checkout) Total(Order) (int, error) { return 0, errors.New("TODO: calculate total") }
func FlatRate(int) ShippingStrategy          { return nil }
func WeightBased(int, int) ShippingStrategy  { return nil }
