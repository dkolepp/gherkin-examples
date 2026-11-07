// The cart team tests its own published kit exactly the way consumers use
// it: by importing it and registering it with godog. The feature files are
// the shared specs in specs/cart/ at the repo root.
package features_test

import (
	"testing"

	"example.com/shop/bddcommon"
	cartbdd "example.com/shop/cart/bdd"
)

func TestCartSpecs(t *testing.T) {
	bddcommon.RunSpecs(t, []string{"../../../../specs/cart"},
		bddcommon.InitializeScenario,
		cartbdd.InitializeScenario,
	)
}
