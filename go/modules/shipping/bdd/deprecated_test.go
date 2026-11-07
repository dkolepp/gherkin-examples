package bdd_test

// The 1.0.0 sentence must keep working in 1.1.x, with a warning.

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"example.com/shop/bddcommon"
	shippingbdd "example.com/shop/shipping/bdd"
)

func TestOldSentenceStillWorksButWarns(t *testing.T) {
	var warnings bytes.Buffer
	saved := bddcommon.DeprecationOutput
	bddcommon.DeprecationOutput = &warnings
	defer func() { bddcommon.DeprecationOutput = saved }()

	ctx := shippingbdd.Install(context.Background())
	shipping := shippingbdd.ShippingFrom(ctx)
	shipping.SetRate("Canada", 1200)
	shipping.Service.Ship("ORD-0001", "Canada", 3)

	if err := shippingbdd.StepShippedOld(ctx, 3, "Canada"); err != nil {
		t.Fatalf("the old sentence should still pass: %v", err)
	}
	if !strings.Contains(warnings.String(), "deprecated since 1.1.0") {
		t.Fatalf("expected a deprecation warning, got %q", warnings.String())
	}
}
