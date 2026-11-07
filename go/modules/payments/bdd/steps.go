package bdd

// Payment steps (owned by the payments team).

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"

	"example.com/shop/bddcommon"
)

type paymentsKey struct{}

// Install puts a fresh PaymentsDriver in ctx. InitializeScenario registers it
// as a Before hook, so each scenario gets its own.
func Install(ctx context.Context) context.Context {
	return context.WithValue(ctx, paymentsKey{}, NewPaymentsDriver())
}

// LookupPayments returns the scenario's PaymentsDriver, and whether there is one.
func LookupPayments(ctx context.Context) (*PaymentsDriver, bool) {
	d, ok := ctx.Value(paymentsKey{}).(*PaymentsDriver)
	return d, ok
}

// PaymentsFrom returns the scenario's PaymentsDriver. It panics, failing the step with a
// clear message, if paymentsbdd.InitializeScenario was not registered.
func PaymentsFrom(ctx context.Context) *PaymentsDriver {
	d, ok := LookupPayments(ctx)
	if !ok {
		panic("no PaymentsDriver in the scenario: register paymentsbdd.InitializeScenario")
	}
	return d
}

// InitializeScenario installs a fresh PaymentsDriver before each scenario and
// registers the payment steps. Register bddcommon.InitializeScenario first.
func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return Install(ctx), nil
	})

	// Give the shopper a card that can be charged up to the limit in total.
	sc.Given(`^the shopper has a card "([^"]*)" with a limit of `+bddcommon.MoneyPattern+`$`, stepCard)
	// Charge a card directly. Fails the scenario if the payment is declined.
	sc.When(`^I charge `+bddcommon.MoneyPattern+` to card "([^"]*)"$`, stepCharge)
	// Try to charge a card, which may be declined.
	// Follow with 'Then I should be told "..."' to check the reason.
	sc.When(`^I try to charge `+bddcommon.MoneyPattern+` to card "([^"]*)"$`, stepTryCharge)
	// Check the total of all successful payments on a card.
	sc.Then(`^card "([^"]*)" should be charged `+bddcommon.MoneyPattern+`$`, stepCharged)
	// Check that no payment went through on a card.
	sc.Then(`^card "([^"]*)" should not be charged$`, stepNotCharged)
}

func stepCard(ctx context.Context, card, limit string) {
	PaymentsFrom(ctx).IssueCard(card, bddcommon.MustMoney(limit))
}

func stepCharge(ctx context.Context, amount, card string) error {
	return PaymentsFrom(ctx).Charge(card, bddcommon.MustMoney(amount))
}

func stepTryCharge(ctx context.Context, amount, card string) {
	bddcommon.OutcomeFrom(ctx).Attempt(func() error {
		return PaymentsFrom(ctx).Charge(card, bddcommon.MustMoney(amount))
	})
}

func stepCharged(ctx context.Context, card, amount string) error {
	if got := PaymentsFrom(ctx).AmountCharged(card); got != bddcommon.MustMoney(amount) {
		return fmt.Errorf("card was charged %s", got)
	}
	return nil
}

func stepNotCharged(ctx context.Context, card string) error {
	if got := PaymentsFrom(ctx).AmountCharged(card); got != 0 {
		return fmt.Errorf("card was charged %s", got)
	}
	return nil
}
