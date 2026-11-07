// The payments team tests its own published kit exactly the way consumers
// use it: by importing it and registering it with godog. The feature files
// are the shared specs in specs/payments/ at the repo root.
package features_test

import (
	"testing"

	"example.com/shop/bddcommon"
	paymentsbdd "example.com/shop/payments/bdd"
)

func TestPaymentsSpecs(t *testing.T) {
	bddcommon.RunSpecs(t, []string{"../../../../specs/payments"},
		bddcommon.InitializeScenario,
		paymentsbdd.InitializeScenario,
	)
}
