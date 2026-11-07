package bdd

// Checkout driver (owned by the checkout team, who own the journey).
//
// It builds the real Checkout service from the other teams' drivers: the cart
// team's cart, the payments team's gateway and the shipping team's service.
// Journey steps call this driver for checkout actions, and the other teams'
// drivers for everything else.

import (
	cartbdd "example.com/shop/cart/bdd"
	"example.com/shop/checkout"
	paymentsbdd "example.com/shop/payments/bdd"
	shippingbdd "example.com/shop/shipping/bdd"
)

// CheckoutDriver drives the checkout application.
type CheckoutDriver struct {
	Shop     *cartbdd.ShopDriver
	Checkout *checkout.Checkout
	Country  string
	Order    *checkout.Order
}

// NewCheckoutDriver builds a checkout from the other teams' drivers.
func NewCheckoutDriver(shop *cartbdd.ShopDriver, payments *paymentsbdd.PaymentsDriver, shipping *shippingbdd.ShippingDriver) *CheckoutDriver {
	return &CheckoutDriver{
		Shop:     shop,
		Checkout: checkout.New(payments.Gateway, shipping.Service),
	}
}

// LiveIn sets the country orders are shipped to.
func (d *CheckoutDriver) LiveIn(country string) { d.Country = country }

// CheckOut places the order for the current cart.
func (d *CheckoutDriver) CheckOut(card string) error {
	order, err := d.Checkout.PlaceOrder(d.Shop.Cart, card, d.Country)
	if err != nil {
		return err
	}
	d.Order = &order
	return nil
}
