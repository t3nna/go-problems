package main

import "fmt"

type Order struct{ SubtotalCents, WeightGrams int }
type ShippingStrategy interface{ ShippingCost(Order) (int, error) }
type Checkout struct{ strategy ShippingStrategy }

func NewCheckout(strategy ShippingStrategy) (*Checkout, error) {
	if strategy == nil {
		return nil, fmt.Errorf("shipping strategy is required")
	}
	return &Checkout{strategy: strategy}, nil
}
func (c *Checkout) Total(order Order) (int, error) {
	if c == nil || c.strategy == nil {
		return 0, fmt.Errorf("shipping strategy is required")
	}
	if order.SubtotalCents < 0 || order.WeightGrams < 0 {
		return 0, fmt.Errorf("order values must not be negative")
	}
	shipping, err := c.strategy.ShippingCost(order)
	if err != nil {
		return 0, err
	}
	return order.SubtotalCents + shipping, nil
}

func FlatRate(cents int) ShippingStrategy {
	return shippingStrategyFunc(func(Order) (int, error) {
		if cents < 0 {
			return 0, fmt.Errorf("flat rate must not be negative")
		}
		return cents, nil
	})
}
func WeightBased(baseCents, perStartedKilogramCents int) ShippingStrategy {
	return shippingStrategyFunc(func(order Order) (int, error) {
		if baseCents < 0 || perStartedKilogramCents < 0 {
			return 0, fmt.Errorf("weight rates must not be negative")
		}
		kilograms := (order.WeightGrams + 999) / 1000
		return baseCents + kilograms*perStartedKilogramCents, nil
	})
}

type shippingStrategyFunc func(Order) (int, error)

func (f shippingStrategyFunc) ShippingCost(order Order) (int, error) { return f(order) }
