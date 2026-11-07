// This journey writes no step code. It registers the teams' published kits
// and runs the shared journey spec, specs/journeys/ at the repo root.
// Order matters: checkout is built from the cart, payments and shipping drivers.
package journey_test

import (
	"testing"

	"example.com/shop/bddcommon"
	cartbdd "example.com/shop/cart/bdd"
	checkoutbdd "example.com/shop/checkout/bdd"
	paymentsbdd "example.com/shop/payments/bdd"
	shippingbdd "example.com/shop/shipping/bdd"
)

func TestJourney(t *testing.T) {
	bddcommon.RunSpecs(t, []string{"../../specs/journeys"},
		bddcommon.InitializeScenario,
		cartbdd.InitializeScenario,
		paymentsbdd.InitializeScenario,
		shippingbdd.InitializeScenario,
		checkoutbdd.InitializeScenario,
	)
}
