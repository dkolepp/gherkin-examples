// Package checkout places orders (owned by the checkout team).
//
// It connects the other domains through their public APIs only:
// cart total -> shipping quote -> card payment -> shipment.
// It is the only application package that imports other teams' apps.
package checkout

import (
	"fmt"

	"example.com/shop/cart"
	"example.com/shop/payments"
	"example.com/shop/shipping"
)

// Error is a checkout refusal of checkout's own.
type Error struct{ Message string }

func (e *Error) Error() string { return e.Message }

// Order is a placed order. Amounts are in cents.
type Order struct {
	ID       string
	Goods    int64
	Shipping int64
	Total    int64
	Payment  payments.Payment
	Shipment shipping.Shipment
}

// Checkout places orders, charging cards and sending shipments.
type Checkout struct {
	Payments *payments.Gateway
	Shipping *shipping.Service
	Orders   []Order
}

// New returns a checkout that uses these payments and shipping services.
func New(p *payments.Gateway, s *shipping.Service) *Checkout {
	return &Checkout{Payments: p, Shipping: s}
}

// PlaceOrder charges the card for the cart plus shipping, ships the items
// and empties the cart. If any part fails, nothing half-happens.
func (c *Checkout) PlaceOrder(basket *cart.Cart, card, country string) (Order, error) {
	if basket.IsEmpty() {
		return Order{}, &Error{"Your cart is empty"}
	}
	goods := basket.Total() // after any discount
	shippingCost, err := c.Shipping.Quote(country, goods)
	if err != nil {
		return Order{}, err // may refuse
	}
	total := goods + shippingCost
	payment, err := c.Payments.Charge(card, total)
	if err != nil {
		return Order{}, err // declined: nothing ships
	}
	orderID := fmt.Sprintf("ORD-%04d", len(c.Orders)+1)
	items := 0
	for _, l := range basket.Lines() {
		items += l.Quantity
	}
	shipment := c.Shipping.Ship(orderID, country, items)
	order := Order{orderID, goods, shippingCost, total, payment, shipment}
	c.Orders = append(c.Orders, order)
	basket.Clear()
	return order, nil
}
