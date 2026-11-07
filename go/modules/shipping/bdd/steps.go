package bdd

// Shipping steps (owned by the shipping team).

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"

	"example.com/shop/bddcommon"
)

type shippingKey struct{}

// Install puts a fresh ShippingDriver in ctx. InitializeScenario registers it
// as a Before hook, so each scenario gets its own.
func Install(ctx context.Context) context.Context {
	return context.WithValue(ctx, shippingKey{}, NewShippingDriver())
}

// LookupShipping returns the scenario's ShippingDriver, and whether there is one.
func LookupShipping(ctx context.Context) (*ShippingDriver, bool) {
	d, ok := ctx.Value(shippingKey{}).(*ShippingDriver)
	return d, ok
}

// ShippingFrom returns the scenario's ShippingDriver. It panics, failing the step with a
// clear message, if shippingbdd.InitializeScenario was not registered.
func ShippingFrom(ctx context.Context) *ShippingDriver {
	d, ok := LookupShipping(ctx)
	if !ok {
		panic("no ShippingDriver in the scenario: register shippingbdd.InitializeScenario")
	}
	return d
}

// InitializeScenario installs a fresh ShippingDriver before each scenario and
// registers the shipping steps. Register bddcommon.InitializeScenario first.
func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return Install(ctx), nil
	})

	// Set a flat shipping rate per country.
	// Columns: country, rate. Countries not listed can't be shipped to.
	sc.Given(`^shipping costs:$`, stepShippingCosts)
	// Make shipping free for orders at or above this value.
	sc.Given(`^orders of `+bddcommon.MoneyPattern+` or more ship free$`, stepFreeFrom)
	// Load shipping rates and the free-shipping threshold from a JSON file.
	// Bundled: "standard_shipping.json". Or give the path of your own file.
	sc.Given(`^shipping is set up from "([^"]*)"$`, stepShippingFromFile)
	// Ask what shipping costs for an order of this value.
	sc.When(`^I get a shipping quote to "([^"]*)" for `+bddcommon.MoneyPattern+`$`, stepQuote)
	// Ask for a shipping quote, which may be refused.
	// Follow with 'Then I should be told "..."' to check the reason.
	sc.When(`^I try to get a shipping quote to "([^"]*)" for `+bddcommon.MoneyPattern+`$`, stepTryQuote)
	// Check the last shipping quote.
	sc.Then(`^shipping should cost `+bddcommon.MoneyPattern+`$`, stepShippingCost)
	// Check that exactly one shipment went out, with this many items.
	sc.Then(`^a shipment of `+bddcommon.QtyPattern+` items should be sent to "([^"]*)"$`, stepShipmentSent)
	// Check that nothing was shipped.
	sc.Then(`^no shipment should be sent$`, stepNoShipment)

	// Version 1.0.0 used this sentence. It was too loose: it didn't say
	// "exactly one shipment". Kept working until 2.0.0 so existing feature
	// files don't break.
	// Deprecated since 1.1.0: use 'a shipment of N items should be sent to "..."'.
	sc.Then(`^`+bddcommon.QtyPattern+` items should be shipped to "([^"]*)"$`, StepShippedOld)
}

func stepShippingCosts(ctx context.Context, table *godog.Table) error {
	for _, row := range bddcommon.Rows(table) {
		rate, err := bddcommon.ToMoney(row["rate"])
		if err != nil {
			return err
		}
		ShippingFrom(ctx).SetRate(row["country"], rate)
	}
	return nil
}

func stepFreeFrom(ctx context.Context, amount string) {
	ShippingFrom(ctx).FreeShippingFrom(bddcommon.MustMoney(amount))
}

func stepShippingFromFile(ctx context.Context, fileName string) error {
	return ShippingFrom(ctx).Load(fileName)
}

func stepQuote(ctx context.Context, country, amount string) error {
	return ShippingFrom(ctx).Quote(country, bddcommon.MustMoney(amount))
}

func stepTryQuote(ctx context.Context, country, amount string) {
	bddcommon.OutcomeFrom(ctx).Attempt(func() error {
		return ShippingFrom(ctx).Quote(country, bddcommon.MustMoney(amount))
	})
}

func stepShippingCost(ctx context.Context, amount string) error {
	quote := ShippingFrom(ctx).LastQuote
	if quote == nil {
		return errors.New("no shipping quote was given")
	}
	if *quote != bddcommon.MustMoney(amount) {
		return fmt.Errorf("quote was %s", *quote)
	}
	return nil
}

func stepShipmentSent(ctx context.Context, qty int, country string) error {
	shipments := ShippingFrom(ctx).Shipments()
	if len(shipments) != 1 {
		return fmt.Errorf("%d shipments were sent", len(shipments))
	}
	if s := shipments[0]; s.Items != qty || s.Country != country {
		return fmt.Errorf("%+v", s)
	}
	return nil
}

func stepNoShipment(ctx context.Context) error {
	if shipments := ShippingFrom(ctx).Shipments(); len(shipments) != 0 {
		return fmt.Errorf("%+v", shipments)
	}
	return nil
}

// StepShippedOld is the 1.0.0 sentence `{qty} items should be shipped to
// "{country}"`. Deprecated since 1.1.0: it warns, then runs the new step.
// It is exported so the unit test can check the warning.
func StepShippedOld(ctx context.Context, qty int, country string) error {
	bddcommon.Deprecated(
		`{qty:Qty} items should be shipped to "{country}"`,
		`a shipment of {qty:Qty} items should be sent to "{country}"`,
		"1.1.0",
	)
	return stepShipmentSent(ctx, qty, country)
}
