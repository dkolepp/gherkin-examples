// The checkout team tests its own published kit exactly the way consumers
// use it: by importing it and registering it with godog. The feature files
// are the shared specs in specs/checkout/ at the repo root.
//
// Order matters: checkout is built from the cart, payments and shipping drivers.
package features_test

import (
	"testing"

	"example.com/shop/bddcommon"
	cartbdd "example.com/shop/cart/bdd"
	checkoutbdd "example.com/shop/checkout/bdd"
	paymentsbdd "example.com/shop/payments/bdd"
	shippingbdd "example.com/shop/shipping/bdd"
)

func TestCheckoutSpecs(t *testing.T) {
	bddcommon.RunSpecs(t, []string{"../../../../specs/checkout"},
		bddcommon.InitializeScenario,
		cartbdd.InitializeScenario,
		paymentsbdd.InitializeScenario,
		shippingbdd.InitializeScenario,
		checkoutbdd.InitializeScenario,
	)
}
