// Plain Go unit tests, for contrast with the BDD scenarios.
//
// BDD scenarios describe behaviour a shopper or product owner cares about.
// Unit tests cover small technical details (rounding, edge cases) that would
// only clutter a .feature file. You usually want both.
package cart_test

import (
	"errors"
	"testing"

	"example.com/shop/cart"
)

func TestToCentsRoundsHalfUpToCents(t *testing.T) {
	for text, want := range map[string]int64{"1.005": 101, "2": 200, "0.004": 0} {
		got, err := cart.ToCents(text)
		if err != nil || got != want {
			t.Errorf("ToCents(%q) = %d, %v; want %d", text, got, err, want)
		}
	}
}

func TestFormatMoney(t *testing.T) {
	if got := cart.FormatMoney(350); got != "$3.50" {
		t.Errorf("got %s", got)
	}
}

func TestRemovingMoreThanInCartJustEmptiesIt(t *testing.T) {
	c := cart.New(cart.NewCatalog(cart.Product{Name: "Mug", Price: 800, Stock: 5}))
	if err := c.Add("Mug", 2); err != nil {
		t.Fatal(err)
	}
	c.Remove("Mug", 10)
	if !c.IsEmpty() {
		t.Fatal("expected the cart to be empty")
	}
}

func TestUnknownProductIsRefused(t *testing.T) {
	err := cart.New(cart.NewCatalog()).Add("Teapot", 1)
	var unknown *cart.UnknownProductError
	if !errors.As(err, &unknown) || err.Error() != "We don't sell Teapot" {
		t.Fatalf("got %v", err)
	}
}
