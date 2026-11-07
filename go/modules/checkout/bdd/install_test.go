package bdd_test

import (
	"context"
	"strings"
	"testing"

	cartbdd "example.com/shop/cart/bdd"
	checkoutbdd "example.com/shop/checkout/bdd"
)

// The contract says checkout needs the other kits' drivers. Registering the
// kits in the wrong order must fail with a message that says what to do.
func TestInstallNamesTheMissingKit(t *testing.T) {
	ctx := cartbdd.Install(context.Background())
	_, err := checkoutbdd.Install(ctx)
	if err == nil || !strings.Contains(err.Error(), "register paymentsbdd.InitializeScenario before") {
		t.Fatalf("got %v", err)
	}
}
