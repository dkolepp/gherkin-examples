# example.com/shop/cart 1.0.0

Owned by the cart team. Covers the cart: products, discount codes and receipts.

This module is two things:

- `example.com/shop/cart`: the application code. It knows nothing about godog.
- `example.com/shop/cart/bdd`: the team's BDD kit, a driver and step library other teams
  can use in their own feature files.

The team's own feature files are the shared, language-neutral specs in
`specs/cart/` at the repo root. `features/` here holds only the godog wiring:
a test that registers the kit, exactly as consumers do, and points godog at
the specs.

## Contract

| | godog |
|---|---|
| Require | `example.com/shop/cart v1.0.0` in your `go.mod` |
| Import | `cartbdd "example.com/shop/cart/bdd"` (every kit's package is named `bdd`) |
| Register | `cartbdd.InitializeScenario`, after `bddcommon.InitializeScenario` |
| Provides | A Before hook, `cartbdd.Install(ctx)`, puts a fresh `*ShopDriver` in the scenario's `context.Context`. Read it with `cartbdd.ShopFrom(ctx)` (panics with a clear message if missing) or `cartbdd.LookupShop(ctx)` |
| Needs | the `*bddcommon.Outcome` |
| Python equivalent | `install(context)` sets `context.shop` |

Bundled test data: `standard_shop.json` (embedded with `go:embed`). A file with that name in your working folder takes precedence.

## Steps

`*` means the step works as Given, When or Then (`sc.Step`); the others are
registered with `sc.Given`, `sc.When` or `sc.Then` and match only that keyword
(an `And` takes the keyword before it).

| | Sentence | What it does |
|---|---|---|
| Given | `the shop sells these products:` | Stock the shop with the products in the table. |
| Given | `these discount codes exist:` | Create the discount codes in the table. |
| Given | `the shop is set up from "{file_name}"` | Stock the shop from a JSON file of products and discount codes. |
| Given | `I have an empty cart` | Start a new, empty cart. |
| Given | `I have a cart containing {qty:Qty} "{product}"` | Start a new cart that already holds some of one product. |
| * | `I add {qty:Qty} "{product}" to the cart` | Put products in the cart. Fails the scenario if that isn't possible. |
| When | `I try to add {qty:Qty} "{product}" to the cart` | Try to put products in the cart, which may be refused. |
| When | `I remove {qty:Qty} "{product}" from the cart` | Take products out of the cart. Removing more than is there empties that line. |
| When | `I apply the discount code "{code}"` | Apply a discount code. Fails the scenario if the code is refused. |
| Given | `I have applied the discount code "{code}"` | Apply a discount code. Fails the scenario if the code is refused. |
| When | `I try to apply the discount code "{code}"` | Try a discount code, which may be refused. |
| When | `I print the receipt` | Print the receipt for the current cart. |
| Then | `the cart should be empty` | Check that the cart holds nothing. |
| Then | `the cart total should be {amount:Money}` | Check the amount to pay, after any discount. Write money as $12.50. |
| Then | `the discount should be {amount:Money}` | Check how much the discount takes off. $0.00 means no discount applies. |
| Then | `the cart should contain:` | Check how many of each product are in the cart. |
| Then | `the cart should contain {qty:Qty} "{product}"` | Check how many of one product are in the cart. |
| Then | `the receipt should read:` | Check the printed receipt matches the text block below the step exactly. |

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
module path in Go (`example.com/shop/cart/v2`): see "Versions" in
[`go/README.md`](../../README.md).

See [CHANGELOG.md](CHANGELOG.md).
