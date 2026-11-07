# example.com/shop/checkout 1.0.0

Owned by the checkout team. Covers placing an order: cart total, shipping quote, card payment, shipment.

This module is two things:

- `example.com/shop/checkout`: the application code. It knows nothing about godog.
  It imports the cart, payments and shipping apps: the only app that imports
  another team's code.
- `example.com/shop/checkout/bdd`: the team's BDD kit, a driver and step library other teams
  can use in their own feature files.

The team's own feature files are the shared, language-neutral specs in
`specs/checkout/` at the repo root. `features/` here holds only the godog wiring:
a test that registers the kit, exactly as consumers do, and points godog at
the specs.

## Contract

| | godog |
|---|---|
| Require | `example.com/shop/checkout v1.0.0` in your `go.mod` |
| Import | `checkoutbdd "example.com/shop/checkout/bdd"` (every kit's package is named `bdd`) |
| Register | `checkoutbdd.InitializeScenario`, after `bddcommon`, `cartbdd`, `paymentsbdd` and `shippingbdd` (last) |
| Provides | A Before hook, `checkoutbdd.Install(ctx)`, puts a fresh `*CheckoutDriver` in the scenario's `context.Context`. Read it with `checkoutbdd.CheckoutFrom(ctx)` (panics with a clear message if missing) or `checkoutbdd.LookupCheckout(ctx)` |
| Needs | the `*cartbdd.ShopDriver`, `*paymentsbdd.PaymentsDriver`, `*shippingbdd.ShippingDriver` and `*bddcommon.Outcome`: register those kits first. If one is missing, the Before hook fails the scenario with a message naming it |
| Python equivalent | `install(context)` sets `context.checkout` |

## Steps

`*` means the step works as Given, When or Then (`sc.Step`); the others are
registered with `sc.Given`, `sc.When` or `sc.Then` and match only that keyword
(an `And` takes the keyword before it).

| | Sentence | What it does |
|---|---|---|
| Given | `the shopper lives in "{country}"` | Set the country the order will be shipped to. |
| * | `I check out paying with card "{card}"` | Place the order for the current cart. Fails the scenario if refused. |
| When | `I try to check out paying with card "{card}"` | Try to place the order, which may be refused by payments or shipping. |
| Then | `the order total should be {amount:Money}` | Check what the order cost in total: goods after discount, plus shipping. |
| Then | `the order number should be "{order_id}"` | Check the number given to the last order placed. |

In Go the placeholders are regular expressions: `{qty:Qty}` is
`bddcommon.QtyPattern` (`(\d+)`), `{amount:Money}` is `bddcommon.MoneyPattern`
(`$12.50`), and `"{name}"` is `"([^"]*)"`. `make -C ../../journey catalog`
prints them.

## Versioning

The step sentences are this module's public API.

- New step: minor version.
- Changed wording or meaning of a step: major version. Deprecate the old sentence
  for at least one minor version first (see `bddcommon.Deprecated`).
- Driver methods and exported Go names follow the same rules.

The version is in [`VERSION`](VERSION); `make build` in `go/` publishes it and
refuses a version with no CHANGELOG entry. A major version also changes the
module path in Go (`example.com/shop/checkout/v2`): see "Versions" in
[`go/README.md`](../../README.md).

See [CHANGELOG.md](CHANGELOG.md).
