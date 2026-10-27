package main

// ShippingStrategy is the part of checkout that may vary.
type ShippingStrategy interface {
	ShippingCost() int
}

// Weight Based
type WeightBased struct {
	Kg         int
	CentsPerKg int
}

func (w WeightBased) ShippingCost() int {
	if w.Kg < 5 {
		return w.Kg * 50
	}
	return w.Kg * 30
}

// FlatRate charges the same amount for every order.
type FlatRate struct {
	Cents int
}

func (r FlatRate) ShippingCost() int { return r.Cents }

// FreeShipping always charges nothing.
type FreeShipping struct{}

func (FreeShipping) ShippingCost() int { return 0 }

type Checkout struct {
	shipping ShippingStrategy
}

func NewCheckout(s ShippingStrategy) Checkout {

	check := Checkout{
		shipping: s,
	}

	return check
}

func (c Checkout) Total(subtototal int) int {
	sum := subtototal + c.shipping.ShippingCost()
	return sum
}
