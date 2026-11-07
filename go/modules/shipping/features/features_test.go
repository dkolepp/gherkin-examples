// The shipping team tests its own published kit exactly the way consumers
// use it: by importing it and registering it with godog. The feature files
// are the shared specs in specs/shipping/ at the repo root.
package features_test

import (
	"testing"

	"example.com/shop/bddcommon"
	shippingbdd "example.com/shop/shipping/bdd"
)

func TestShippingSpecs(t *testing.T) {
	bddcommon.RunSpecs(t, []string{"../../../../specs/shipping"},
		bddcommon.InitializeScenario,
		shippingbdd.InitializeScenario,
	)
}
