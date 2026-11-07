package bdd

// Payments driver (owned by the payments team). Part of the public contract.
//
// Other teams' steps use these methods instead of touching the Gateway
// directly, so the payments team can change how payments work inside.

import (
	"example.com/shop/bddcommon"
	"example.com/shop/payments"
)

// PaymentsDriver drives the payments application.
type PaymentsDriver struct {
	Gateway *payments.Gateway
}

// NewPaymentsDriver returns a driver for a gateway with no cards.
func NewPaymentsDriver() *PaymentsDriver {
	return &PaymentsDriver{Gateway: payments.NewGateway()}
}

// IssueCard gives the shopper a card with a spending limit.
func (d *PaymentsDriver) IssueCard(number string, limit bddcommon.Money) {
	d.Gateway.AddCard(number, int64(limit))
}

// Charge charges a card.
func (d *PaymentsDriver) Charge(number string, amount bddcommon.Money) error {
	_, err := d.Gateway.Charge(number, int64(amount))
	return err
}

// AmountCharged is the total of all successful payments on a card.
func (d *PaymentsDriver) AmountCharged(number string) bddcommon.Money {
	return bddcommon.Money(d.Gateway.ChargedTo(number))
}
