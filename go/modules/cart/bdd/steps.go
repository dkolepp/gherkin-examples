package bdd

// Cart steps (owned by the cart team).
//
// Each function matches one Gherkin sentence, pulls the values out of it, and
// hands them to the ShopDriver. No business logic lives here.

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cucumber/godog"

	"example.com/shop/bddcommon"
)

type shopKey struct{}

// Install puts a fresh ShopDriver in ctx. InitializeScenario registers it as
// a Before hook, so each scenario gets its own.
func Install(ctx context.Context) context.Context {
	return context.WithValue(ctx, shopKey{}, NewShopDriver())
}

// LookupShop returns the scenario's ShopDriver, and whether there is one.
func LookupShop(ctx context.Context) (*ShopDriver, bool) {
	d, ok := ctx.Value(shopKey{}).(*ShopDriver)
	return d, ok
}

// ShopFrom returns the scenario's ShopDriver. It panics, failing the step with a
// clear message, if cartbdd.InitializeScenario was not registered.
func ShopFrom(ctx context.Context) *ShopDriver {
	d, ok := LookupShop(ctx)
	if !ok {
		panic("no ShopDriver in the scenario: register cartbdd.InitializeScenario")
	}
	return d
}

// InitializeScenario installs a fresh ShopDriver before each scenario and
// registers the cart steps. Register bddcommon.InitializeScenario first.
func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return Install(ctx), nil
	})

	// --- Given: put the world in a known state ---

	// Stock the shop with the products in the table.
	// Columns: name, price, stock. A stock of 0 means sold out.
	sc.Given(`^the shop sells these products:$`, stepShopSells)
	// Create the discount codes in the table.
	// Columns: code, percent off, minimum spend (use 0.00 for no minimum).
	sc.Given(`^these discount codes exist:$`, stepDiscountCodes)
	// Stock the shop from a JSON file of products and discount codes.
	// Bundled: "standard_shop.json". Or give the path of your own file.
	sc.Given(`^the shop is set up from "([^"]*)"$`, stepShopFromFile)
	// Start a new, empty cart.
	sc.Given(`^I have an empty cart$`, stepEmptyCart)
	// Start a new cart that already holds some of one product.
	sc.Given(`^I have a cart containing `+bddcommon.QtyPattern+` "([^"]*)"$`, stepCartContaining)

	// --- When: the action under test ---
	// sc.Step matches Given, When and Then, so one sentence can be used for
	// setup ("And I add ...") as well as for the action itself.

	// Put products in the cart. Fails the scenario if that isn't possible.
	sc.Step(`^I add `+bddcommon.QtyPattern+` "([^"]*)" to the cart$`, stepAdd)
	// Try to put products in the cart, which may be refused.
	// Follow with 'Then I should be told "..."' to check the reason.
	sc.When(`^I try to add `+bddcommon.QtyPattern+` "([^"]*)" to the cart$`, stepTryAdd)
	// Take products out of the cart. Removing more than is there empties that line.
	sc.When(`^I remove `+bddcommon.QtyPattern+` "([^"]*)" from the cart$`, stepRemove)
	// Apply a discount code. Fails the scenario if the code is refused.
	// Two sentences, one function.
	sc.When(`^I apply the discount code "([^"]*)"$`, stepApplyCode)
	sc.Given(`^I have applied the discount code "([^"]*)"$`, stepApplyCode)
	// Try a discount code, which may be refused.
	// Follow with 'Then I should be told "..."' to check the reason.
	sc.When(`^I try to apply the discount code "([^"]*)"$`, stepTryApplyCode)
	// Print the receipt for the current cart.
	sc.When(`^I print the receipt$`, stepPrintReceipt)

	// --- Then: check the outcome ---

	// Check that the cart holds nothing.
	sc.Then(`^the cart should be empty$`, stepCartEmpty)
	// Check the amount to pay, after any discount. Write money as $12.50.
	sc.Then(`^the cart total should be `+bddcommon.MoneyPattern+`$`, stepCartTotal)
	// Check how much the discount takes off. $0.00 means no discount applies.
	sc.Then(`^the discount should be `+bddcommon.MoneyPattern+`$`, stepDiscount)
	// Check how many of each product are in the cart. Columns: product, quantity.
	sc.Then(`^the cart should contain:$`, stepCartContainsTable)
	// Check how many of one product are in the cart.
	sc.Then(`^the cart should contain `+bddcommon.QtyPattern+` "([^"]*)"$`, stepCartContains)
	// Check the printed receipt matches the text block below the step exactly.
	sc.Then(`^the receipt should read:$`, stepReceipt)
}

// atoi converts a table cell to an int.
func atoi(cell string) (int, error) { return strconv.Atoi(cell) }

func stepShopSells(ctx context.Context, table *godog.Table) error {
	for _, row := range bddcommon.Rows(table) {
		stock, err := atoi(row["stock"])
		if err != nil {
			return err
		}
		if err := ShopFrom(ctx).StockProduct(row["name"], row["price"], stock); err != nil {
			return err
		}
	}
	return nil
}

func stepDiscountCodes(ctx context.Context, table *godog.Table) error {
	for _, row := range bddcommon.Rows(table) {
		percent, err := atoi(row["percent off"])
		if err != nil {
			return err
		}
		if err := ShopFrom(ctx).AddDiscountCode(row["code"], percent, row["minimum spend"]); err != nil {
			return err
		}
	}
	return nil
}

func stepShopFromFile(ctx context.Context, fileName string) error {
	return ShopFrom(ctx).LoadShop(fileName)
}

func stepEmptyCart(ctx context.Context) {
	ShopFrom(ctx).NewCart()
}

func stepCartContaining(ctx context.Context, qty int, product string) error {
	shop := ShopFrom(ctx)
	shop.NewCart()
	return shop.AddToCart(product, qty)
}

func stepAdd(ctx context.Context, qty int, product string) error {
	return ShopFrom(ctx).AddToCart(product, qty)
}

func stepTryAdd(ctx context.Context, qty int, product string) {
	bddcommon.OutcomeFrom(ctx).Attempt(func() error { return ShopFrom(ctx).AddToCart(product, qty) })
}

func stepRemove(ctx context.Context, qty int, product string) {
	ShopFrom(ctx).RemoveFromCart(product, qty)
}

func stepApplyCode(ctx context.Context, code string) error {
	return ShopFrom(ctx).ApplyCode(code)
}

func stepTryApplyCode(ctx context.Context, code string) {
	bddcommon.OutcomeFrom(ctx).Attempt(func() error { return ShopFrom(ctx).ApplyCode(code) })
}

func stepPrintReceipt(ctx context.Context) {
	ShopFrom(ctx).PrintReceipt()
}

func stepCartEmpty(ctx context.Context) error {
	if !ShopFrom(ctx).Cart.IsEmpty() {
		return errors.New("expected the cart to be empty")
	}
	return nil
}

func stepCartTotal(ctx context.Context, amount string) error {
	want := bddcommon.MustMoney(amount)
	if got := bddcommon.Money(ShopFrom(ctx).Cart.Total()); got != want {
		return fmt.Errorf("total was %s", got)
	}
	return nil
}

func stepDiscount(ctx context.Context, amount string) error {
	want := bddcommon.MustMoney(amount)
	if got := bddcommon.Money(ShopFrom(ctx).Cart.DiscountAmount()); got != want {
		return fmt.Errorf("discount was %s", got)
	}
	return nil
}

func stepCartContainsTable(ctx context.Context, table *godog.Table) error {
	for _, row := range bddcommon.Rows(table) {
		want, err := atoi(row["quantity"])
		if err != nil {
			return err
		}
		if got := ShopFrom(ctx).Cart.QuantityOf(row["product"]); got != want {
			return fmt.Errorf("%s: %d", row["product"], got)
		}
	}
	return nil
}

func stepCartContains(ctx context.Context, qty int, product string) error {
	if got := ShopFrom(ctx).Cart.QuantityOf(product); got != qty {
		return fmt.Errorf("cart has %d %s", got, product)
	}
	return nil
}

func stepReceipt(ctx context.Context, text *godog.DocString) error {
	// text.Content holds the doc string (the text between the """ marks).
	if got := ShopFrom(ctx).Receipt; got != text.Content {
		return fmt.Errorf("receipt was:\n%s", got)
	}
	return nil
}
