package bdd

// Checkout steps (owned by the checkout team).

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"

	"example.com/shop/bddcommon"
	cartbdd "example.com/shop/cart/bdd"
	paymentsbdd "example.com/shop/payments/bdd"
	shippingbdd "example.com/shop/shipping/bdd"
)

type checkoutKey struct{}

// Install puts a CheckoutDriver in ctx, built from the ShopDriver,
// PaymentsDriver and ShippingDriver already there. Those kits' hooks must
// have run first; if one hasn't, Install returns an error naming it.
func Install(ctx context.Context) (context.Context, error) {
	shop, okShop := cartbdd.LookupShop(ctx)
	payments, okPayments := paymentsbdd.LookupPayments(ctx)
	shipping, okShipping := shippingbdd.LookupShipping(ctx)
	switch {
	case !okShop:
		return ctx, errors.New(needs("cartbdd"))
	case !okPayments:
		return ctx, errors.New(needs("paymentsbdd"))
	case !okShipping:
		return ctx, errors.New(needs("shippingbdd"))
	}
	return context.WithValue(ctx, checkoutKey{}, NewCheckoutDriver(shop, payments, shipping)), nil
}

func needs(kit string) string {
	return "the checkout kit needs " + kit + ": register " + kit + ".InitializeScenario before checkoutbdd.InitializeScenario"
}

// LookupCheckout returns the scenario's CheckoutDriver, and whether there is one.
func LookupCheckout(ctx context.Context) (*CheckoutDriver, bool) {
	d, ok := ctx.Value(checkoutKey{}).(*CheckoutDriver)
	return d, ok
}

// CheckoutFrom returns the scenario's CheckoutDriver. It panics, failing the
// step with a clear message, if checkoutbdd.InitializeScenario was not registered.
func CheckoutFrom(ctx context.Context) *CheckoutDriver {
	d, ok := LookupCheckout(ctx)
	if !ok {
		panic("no CheckoutDriver in the scenario: register checkoutbdd.InitializeScenario")
	}
	return d
}

// InitializeScenario installs a CheckoutDriver before each scenario and
// registers the checkout steps. Register it after bddcommon, cartbdd,
// paymentsbdd and shippingbdd: its Before hook needs their drivers.
func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return Install(ctx)
	})

	// Set the country the order will be shipped to.
	sc.Given(`^the shopper lives in "([^"]*)"$`, stepLivesIn)
	// Place the order for the current cart. Fails the scenario if refused.
	// Works as Given, When or Then.
	sc.Step(`^I check out paying with card "([^"]*)"$`, stepCheckOut)
	// Try to place the order, which may be refused by payments or shipping.
	// Follow with 'Then I should be told "..."' to check the reason.
	sc.When(`^I try to check out paying with card "([^"]*)"$`, stepTryCheckOut)
	// Check what the order cost in total: goods after discount, plus shipping.
	sc.Then(`^the order total should be `+bddcommon.MoneyPattern+`$`, stepOrderTotal)
	// Check the number given to the last order placed.
	sc.Then(`^the order number should be "([^"]*)"$`, stepOrderNumber)
}

func stepLivesIn(ctx context.Context, country string) {
	CheckoutFrom(ctx).LiveIn(country)
}

func stepCheckOut(ctx context.Context, card string) error {
	return CheckoutFrom(ctx).CheckOut(card)
}

func stepTryCheckOut(ctx context.Context, card string) {
	bddcommon.OutcomeFrom(ctx).Attempt(func() error { return CheckoutFrom(ctx).CheckOut(card) })
}

func stepOrderTotal(ctx context.Context, amount string) error {
	order := CheckoutFrom(ctx).Order
	if order == nil {
		return errors.New("no order was placed")
	}
	if got := bddcommon.Money(order.Total); got != bddcommon.MustMoney(amount) {
		return fmt.Errorf("order total was %s", got)
	}
	return nil
}

func stepOrderNumber(ctx context.Context, orderID string) error {
	order := CheckoutFrom(ctx).Order
	if order == nil {
		return errors.New("no order was placed")
	}
	if order.ID != orderID {
		return fmt.Errorf("order number was %s", order.ID)
	}
	return nil
}
