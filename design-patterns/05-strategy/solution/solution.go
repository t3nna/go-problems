package main

// ShippingStrategy is the part of checkout that may vary.
type ShippingStrategy interface {
	ShippingCost() int
}

// FlatRate charges the same amount for every order.
type FlatRate struct {
	Cents int
}

func (rate FlatRate) ShippingCost() int { return rate.Cents }

// FreeShipping always charges nothing.
type FreeShipping struct{}

func (FreeShipping) ShippingCost() int { return 0 }

type Checkout struct {
	shipping ShippingStrategy
}

func NewCheckout(strategy ShippingStrategy) Checkout {
	return Checkout{shipping: strategy}
}

func (checkout Checkout) Total(subtotalCents int) int {
	return subtotalCents + checkout.shipping.ShippingCost()
}
