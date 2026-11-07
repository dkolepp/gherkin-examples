package bddcommon

import (
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"
)

// Contract
//
//	Provides  Install(ctx) puts a fresh *Outcome in the scenario's context;
//	          OutcomeFrom(ctx) reads it back.
//	Needs     nothing. Register InitializeScenario before every other kit.

type outcomeKey struct{}

// Install puts a fresh Outcome in ctx. InitializeScenario registers it as a
// Before hook, so each scenario gets its own.
func Install(ctx context.Context) context.Context {
	return context.WithValue(ctx, outcomeKey{}, &Outcome{})
}

// LookupOutcome returns the scenario's Outcome, and whether there is one.
func LookupOutcome(ctx context.Context) (*Outcome, bool) {
	o, ok := ctx.Value(outcomeKey{}).(*Outcome)
	return o, ok
}

// OutcomeFrom returns the scenario's Outcome. It panics, failing the step
// with a clear message, if bddcommon.InitializeScenario was not registered.
func OutcomeFrom(ctx context.Context) *Outcome {
	o, ok := LookupOutcome(ctx)
	if !ok {
		panic("no Outcome in the scenario: register bddcommon.InitializeScenario before the other kits")
	}
	return o
}

// InitializeScenario installs a fresh Outcome before each scenario and
// registers the shared step.
func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return Install(ctx), nil
	})

	// Check the message the shopper saw when the last 'try to' step was refused.
	sc.Then(`^I should be told "([^"]*)"$`, stepToldThat)
}

func stepToldThat(ctx context.Context, message string) error {
	err := OutcomeFrom(ctx).LastError
	if err == nil {
		return errors.New("expected an error, but the action succeeded")
	}
	if err.Error() != message {
		return fmt.Errorf("got: %s", err)
	}
	return nil
}
